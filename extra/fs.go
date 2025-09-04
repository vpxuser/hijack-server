package extra

import (
	"github.com/vpxuser/proxy"
	"net/http"
)

func FileServer() {
	proxy.Info("文件服务地址：http://0.0.0.0:8081")
	fs := http.FileServer(http.Dir("./report"))
	http.Handle("/", fs)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		proxy.Fatal(err)
	}
}
