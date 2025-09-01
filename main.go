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
	go hijackServer()
	go extra.CaptureServer()
	extra.FileServer()
}

func hijackServer() {
	cfg := proxy.NewConfig(proxy.From(setting.Cert, setting.Key))
	cfg.Dispatcher = hook.ConnectHandler
	cfg.ClientTLSConfig.InsecureSkipVerify = true
	cfg.WithReqMatcher().Handle(handler.Request)
	cfg.WithRespMatcher().Handle(handler.Response)
	addr := net.JoinHostPort(setting.Cfg.Host, setting.Cfg.Port)
	err := proxy.ListenAndServe(addr, cfg)
	if err != nil {
		proxy.Fatal(err)
	}
}
