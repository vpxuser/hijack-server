package capture

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"github.com/vpxuser/proxy"
	"hijack-server/extra/capture/handler"
	"hijack-server/hook"
	"hijack-server/setting"
	"hijack-server/tools"
	"strings"
)

func CaptureProxy() {
	proxy.Infof("抓包代理服务地址：http://%s", strings.ReplaceAll(setting.Cfg.CaptureProxy, "0.0.0.0", "127.0.0.1"))
	var (
		cert *x509.Certificate
		key  crypto.PrivateKey
	)
	tools.LoadCert("config/cacert.pem", &cert)
	tools.LoadKey("config/cakey.pem", &key)
	cfg := proxy.NewConfig(proxy.FromCA(cert, key))
	cfg.ClientTLSConfig = &tls.Config{InsecureSkipVerify: true}
	cfg.WithReqMatcher().Handle(handler.Req)
	cfg.WithRespMatcher().Handle(handler.Resp)
	if err := proxy.ListenAndServe(setting.Cfg.CaptureProxy, cfg); err != nil {
		proxy.Fatal(err)
	}
}

func CaptureWifi() {
	proxy.Infof("抓包热点服务地址：%s", setting.Cfg.CaptureWifi)
	var (
		cert *x509.Certificate
		key  crypto.PrivateKey
	)
	tools.LoadCert("config/cacert.pem", &cert)
	tools.LoadKey("config/cakey.pem", &key)
	cfg := proxy.NewConfig(proxy.FromCA(cert, key))
	cfg.Negotiator = hook.TProxy
	cfg.ClientTLSConfig = &tls.Config{InsecureSkipVerify: true}
	cfg.WithReqMatcher().Handle(handler.Req)
	cfg.WithRespMatcher().Handle(handler.Resp)
	if err := proxy.ListenAndServe(setting.Cfg.CaptureWifi, cfg); err != nil {
		proxy.Fatal(err)
	}
}
