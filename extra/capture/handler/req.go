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
	"sync"
)

var hosts = new(sync.Map)

const record = false

func Req(r *http.Request, ctx *proxy.Context) (*http.Request, *http.Response) {
	if record {
		if _, ok := hosts.Load(r.Host); !ok {
			hosts.Store(r.Host, nil)
		}
		sb := new(strings.Builder)
		hosts.Range(func(host, none interface{}) bool {
			sb.WriteString(host.(string) + "\n")
			return true
		})
		ctx.Infof("\n%s", sb.String())
	}

	r.Header.Set("X-Request-ID", uuid.New().String())

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		body, err := io.ReadAll(r.Body)
		if err == nil && len(body) > 0 {
			reset := true
			defer func() {
				if reset {
					r.ContentLength = int64(len(body))
					r.Body = io.NopCloser(bytes.NewReader(body))
				}
			}()

			pettyBody, err := json.MarshalIndent(body, "", "  ")
			if err != nil {
				reset = false
				r.Body = io.NopCloser(bytes.NewReader(body))
			} else {
				r.ContentLength = int64(len(pettyBody))
				r.Body = io.NopCloser(bytes.NewReader(pettyBody))
			}
		}
	}

	snapshot, _ := httputil.DumpRequest(r, true)
	ctx.Infof("是否为 https：%v，请求报文内容：\n%s", ctx.Conn.IsTLS(), snapshot)
	return r, nil
}
