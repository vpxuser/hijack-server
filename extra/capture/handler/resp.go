package handler

import (
	"bytes"
	"encoding/json"
	"github.com/vpxuser/proxy"
	"hijack-server/tools"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

func Resp(r *http.Response, ctx *proxy.Context) *http.Response {
	r.Header.Set("X-Request-ID", r.Request.Header.Get("X-Request-ID"))
	body, err := io.ReadAll(r.Body)
	if err == nil && len(body) > 0 {
		defer func() {
			r.ContentLength = int64(len(body))
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}()

		gzipBody := body
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") &&
			body[0] == 0x1f &&
			body[1] == 0x8b {
			gzipBody, err = tools.GzipDecompress(body)
			if err != nil {
				ctx.Error(err)
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				return r
			}
		}

		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			pettyBody, err := json.MarshalIndent(gzipBody, "", "  ")
			if err != nil {
				ctx.Warn(err)
				r.ContentLength = int64(len(gzipBody))
				r.Body = io.NopCloser(bytes.NewBuffer(gzipBody))
			} else {
				r.ContentLength = int64(len(pettyBody))
				r.Body = io.NopCloser(bytes.NewBuffer(pettyBody))
			}
		}
	}

	dump, _ := httputil.DumpResponse(r, true)
	ctx.Infof("响应报文内容：\n%s", dump)
	return r
}
