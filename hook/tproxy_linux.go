//go:build linux
// +build linux

package hook

import (
	"fmt"
	"github.com/vpxuser/proxy"
	"golang.org/x/sys/unix"
	"net"
	"strconv"
	"syscall"
	"unsafe"
)

type address struct {
	typ  uint16
	port [2]byte
	host [4]byte
	zero [8]byte
}

var TProxy proxy.HandshakeFn = func(ctx *proxy.Context) error {
	fd, err := ctx.Conn.Conn.(*net.TCPConn).File()
	if err != nil {
		return err
	}
	defer fd.Close()

	var (
		addr   = address{}
		vallen = uint32(unsafe.Sizeof(addr))
	)

	_, _, errno := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		fd.Fd(),
		uintptr(syscall.SOL_IP),
		uintptr(unix.SO_ORIGINAL_DST),
		uintptr(unsafe.Pointer(&addr)),
		uintptr(unsafe.Pointer(&vallen)),
		0,
	)
	if errno != 0 {
		return err
	}

	ctx.DstHost = fmt.Sprintf("%d.%d.%d.%d", addr.host[0], addr.host[1], addr.host[2], addr.host[3])
	ctx.DstPort = strconv.Itoa(int(uint16(addr.port[0])<<8 + uint16(addr.port[1])))
	ctx.Debugf("解析源目标地址为：%s", net.JoinHostPort(ctx.DstHost, ctx.DstPort))
	return nil
}
