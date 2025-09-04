package main

import (
	"github.com/vpxuser/proxy"
	"hijack-server/extra"
	"hijack-server/handler"
	"hijack-server/hook"
	"hijack-server/setting"
	"net"
)

func init() {
	proxy.SetLogLevel(setting.Cfg.LogLevel)
}

func main() {
	go extra.CaptureServer()
	go extra.FileServer()
	hijackServer()
}

func hijackServer() {
	cfg := proxy.NewConfig(proxy.From(setting.Cert, setting.Key))
	if setting.Cfg.Hotspot {
		cfg.Negotiator = hook.TProxy
	}
	cfg.Dispatcher = hook.ConnectHandler
	cfg.ClientTLSConfig.InsecureSkipVerify = true
	cfg.WithReqMatcher().Handle(handler.Request)
	cfg.WithRespMatcher().Handle(handler.Response)
	addr := net.JoinHostPort(setting.Cfg.Host, setting.Cfg.Port)
	proxy.Infof("劫持服务地址：%s", addr)
	err := proxy.ListenAndServe(addr, cfg)
	if err != nil {
		proxy.Fatal(err)
	}
}
