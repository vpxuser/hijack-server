package hook

import (
	"bufio"
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/inconshreveable/go-vhost"
	"github.com/vpxuser/proxy"
	"hijack-server/setting"
	"net"
	"net/http"
	"strings"
	"sync"
)

var blackList = new(sync.Map)

var ConnectHandler proxy.DispatchFn = func(ctx *proxy.Context) error {
	req, parseErr := http.ReadRequest(bufio.NewReader(ctx.Conn.PeekRd)) //预读取数据，查看是否能解析为http请求
	if parseErr != nil {                                                //解析失败，说明数据流不是http明文请求
		raw, err := ctx.Conn.Peek(2) //预读取前两个字节
		if len(raw) <= 0 {           //没有数据，不符合预期，直接结束进程
			ctx.Error(err)
			return err
		}

		if raw[0] == 0x16 && raw[1] == 0x03 { //识别数据流为tls协议
			//获取SNI
			serverName := ctx.DstHost
			if !proxy.IsDomain(serverName) {
				rawConn, err := vhost.TLS(ctx.Conn)
				if err != nil {
					return ctx.TcpHandler.HandleTcp(ctx)
				}
				serverName = rawConn.Host()
				ctx.Conn = proxy.NewConn(rawConn)
			}

			//查看SNI是否在劫持名单内
			if _, ok := setting.Targets[serverName]; !ok {
				return ctx.TcpHandler.HandleTcp(ctx)
			}

			if setting.Cfg.Skip {
				proxyAddr := net.JoinHostPort(ctx.DstHost, ctx.DstPort)
				//查看缓存，哪些SNI无法进行劫持，直接跳过
				_, ok := blackList.Load(proxyAddr)
				if ok {
					ctx.Infof("跳过 %s tls 劫持", proxyAddr)
					return ctx.TcpHandler.HandleTcp(ctx)
				}
			}

			//返回伪造证书给客户端
			tlsCfg, err := ctx.TLSConfig.From(serverName)
			if err != nil {
				ctx.Error(err)
				return err
			}
			ctx.Conn = proxy.NewConn(tls.Server(ctx.Conn, tlsCfg))
		}

		//预读取请求
		req, err = http.ReadRequest(bufio.NewReader(ctx.Conn.PeekRd))
		if err != nil {
			if setting.Cfg.Skip && strings.Contains(err.Error(), "tls: unknown certificate") { //劫持失败，将此SNI加入黑名单，下次直接跳过劫持
				proxyAddr := net.JoinHostPort(ctx.DstHost, ctx.DstPort)
				ctx.Warnf("对 %s 进行 tls 劫持失败", proxyAddr)
				blackList.Store(proxyAddr, struct{}{})
				return err
			}
			return ctx.TcpHandler.HandleTcp(ctx)
		}
	}

	//由于http明文没有进入tls协议逻辑，所以需要再次进行名单过滤
	//根据请求的Host头，再次判断目标主机是否在劫持名单里
	if ctx.Conn.IsTLS() {
		hostname, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			hostname = req.Host
		}

		if _, ok := setting.Targets[hostname]; !ok {
			return ctx.TcpHandler.HandleTcp(ctx)
		}
	}

	ctx.Debugf("开始劫持 %s:%s", ctx.DstHost, ctx.DstPort)

	if websocket.IsWebSocketUpgrade(req) { //如果是websocket协议，进入ws劫持
		return ctx.WsHandler.HandleWs(ctx)
	}

	return ctx.HttpHandler.HandleHttp(ctx) //进入http劫持
}
