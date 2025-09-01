package handler

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/vpxuser/proxy"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

func Request(req *http.Request, ctx *proxy.Context) (*http.Request, *http.Response) {
	//检查协议
	if req.URL.Scheme == "" {
		if ctx.Conn.IsTLS() {
			req.URL.Scheme = "https"
		} else {
			req.URL.Scheme = "http"
		}
	}

	//检查目标地址
	if req.URL.Hostname() == "" {
		req.URL.Host = req.Host
	}

	//填充响应的唯一id，方便日志追踪
	req.Header.Set("X-Request-ID", uuid.New().String())

	//格式化json请求体
	if strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		body, err := io.ReadAll(req.Body)
		if err == nil {
			pettyBody, err := json.MarshalIndent(json.RawMessage(body), "", "  ")
			if err == nil {
				req.ContentLength = int64(len(pettyBody))
				req.Body = io.NopCloser(bytes.NewReader(pettyBody))
			}
		}
	}

	//打印请求
	snapshot, _ := httputil.DumpRequestOut(req, true)
	ctx.Debugf("协议：%s，请求：\n%s", req.URL.Scheme, snapshot)
	ctx.Extra = string(snapshot)
	return req, nil
}
