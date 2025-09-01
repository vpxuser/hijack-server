package extra

import (
	"github.com/vpxuser/proxy"
	"net/http"
)

func FileServer() {
	fs := http.FileServer(http.Dir("./report"))
	http.Handle("/", fs)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		proxy.Fatal(err)
	}
}
