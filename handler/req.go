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

func Request(r *http.Request, ctx *proxy.Context) (*http.Request, *http.Response) {
	//检查协议
	if r.URL.Scheme == "" {
		if ctx.Conn.IsTLS() {
			r.URL.Scheme = "https"
		} else {
			r.URL.Scheme = "http"
		}
	}

	//检查目标地址
	if r.URL.Hostname() == "" {
		r.URL.Host = r.Host
	}

	//填充响应的唯一id，方便日志追踪
	r.Header.Set("X-Request-ID", uuid.New().String())

	//格式化json请求体
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		body, err := io.ReadAll(r.Body)
		if err == nil && len(body) > 0 {
			reset := true
			defer func() {
				if reset {
					r.ContentLength = int64(len(body))
					r.Body = io.NopCloser(bytes.NewBuffer(body))
				}
			}()

			pettyBody, err := json.MarshalIndent(json.RawMessage(body), "", "  ")
			if err != nil {
				reset = false
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			} else {
				r.ContentLength = int64(len(pettyBody))
				r.Body = io.NopCloser(bytes.NewReader(pettyBody))
			}
		}
	}

	//打印请求
	snapshot, _ := httputil.DumpRequestOut(r, true)
	ctx.Debugf("协议：%s，请求：\n%s", r.URL.Scheme, snapshot)
	ctx.Extra = string(snapshot)
	return r, nil
}
