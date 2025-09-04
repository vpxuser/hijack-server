package extra

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/vpxuser/proxy"
	"hijack-server/tools"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

func CaptureServer() {
	proxy.Info("抓包服务地址：0.0.0.0:8080")
	var (
		cert *x509.Certificate
		key  crypto.PrivateKey
	)

	tools.LoadCert("config/cacrt.pem", &cert)
	tools.LoadKey("config/cakey.pem", &key)

	cfg := proxy.NewConfig(proxy.FromCA(cert, key))

	cfg.ClientTLSConfig = &tls.Config{InsecureSkipVerify: true}

	hosts := new(sync.Map)

	cfg.WithReqMatcher().Handle(func(req *http.Request, ctx *proxy.Context) (*http.Request, *http.Response) {
		if _, ok := hosts.Load(req.Host); !ok {
			hosts.Store(req.Host, nil)
		}
		sb := new(strings.Builder)
		hosts.Range(func(host, none interface{}) bool {
			sb.WriteString(host.(string) + "\n")
			return true
		})
		ctx.Infof("\n%s", sb.String())
		req.Header.Set("X-Request-ID", uuid.New().String())
		dump, _ := httputil.DumpRequest(req, true)
		ctx.Debugf("是否为 https：%v，请求报文内容：\n%s", ctx.Conn.IsTLS(), dump)
		return req, nil
	})

	cfg.WithRespMatcher().Handle(func(resp *http.Response, ctx *proxy.Context) *http.Response {
		resp.Header.Set("X-Request-ID", resp.Request.Header.Get("X-Request-ID"))
		compressed := false
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.Warn(err)
		} else {
			if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") &&
				body[0] == 0x1f &&
				body[1] == 0x8b {
				compressed = true
				body, err = tools.GzipDecompress(body)
				if err != nil {
					ctx.Error(err)
					resp.Body = io.NopCloser(bytes.NewBuffer(body))
					return resp
				}
			}

			if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
				pettyBody, _ := json.MarshalIndent(json.RawMessage(body), "", "  ")
				resp.ContentLength = int64(len(pettyBody))
				resp.Body = io.NopCloser(bytes.NewBuffer(pettyBody))
			} else {
				resp.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		dump, _ := httputil.DumpResponse(resp, true)
		ctx.Debugf("响应报文内容：\n%s", dump)

		if compressed {
			body, err = tools.GzipCompress(body)
			if err != nil {
				ctx.Error(err)
				resp.Body = io.NopCloser(bytes.NewBuffer(body))
				return resp
			}
			resp.ContentLength = int64(len(body))
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		return resp
	})

	if err := proxy.ListenAndServe("0.0.0.0:8080", cfg); err != nil {
		proxy.Fatal(err)
	}
}
