package extra

import (
	"github.com/vpxuser/proxy"
	"hijack-server/setting"
	"net/http"
	"strings"
)

func FileServer() {
	proxy.Infof("文件服务地址：http://%s", strings.ReplaceAll(setting.Cfg.FileServerHost, "0.0.0.0", "127.0.0.1"))
	fs := http.FileServer(http.Dir("./report"))
	http.Handle("/", fs)
	err := http.ListenAndServe(setting.Cfg.FileServerHost, nil)
	if err != nil {
		proxy.Fatal(err)
	}
}
