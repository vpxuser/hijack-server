//go:build windows
// +build windows

package hook

import (
	"github.com/vpxuser/proxy"
)

var TProxy proxy.HandshakeFn = func(ctx *proxy.Context) error {
	return proxy.HttpNegotiator.Handshake(ctx)
}
