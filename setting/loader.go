package setting

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/vpxuser/proxy"
	"hijack-server/tools"
	"os"
	"regexp"
	"strings"
)

const Main_Config_Path = "config/config.yml"

var (
	Cfg      = new(Config)
	Cert     *x509.Certificate
	Key      crypto.PrivateKey
	Targets  = make(map[string]struct{})
	Headers  = make(map[string][]string)
	QRCode   = make(map[string]*regexp.Regexp)
	Image    []byte
	Html     []byte
	Template *os.File
)

func init() {
	loadCfg(Main_Config_Path, Cfg)
	proxy.Debugf("主配置 %s 加载成功", Main_Config_Path)

	tools.LoadCert(Cfg.Cert, &Cert)
	proxy.Debugf("劫持证书 %s 加载成功", Cfg.Cert)

	sb := new(strings.Builder)
	sb.WriteString("证书信息摘要：\n")
	sb.WriteString(fmt.Sprintf("       ├── 证书颁发机构：%s\n", Cert.Issuer.String()))
	sb.WriteString(fmt.Sprintf("       ├── 证书域名：%s\n", Cert.Subject.String()))
	sb.WriteString(fmt.Sprintf("       └── 证书失效日期：%s", Cert.NotAfter.String()))
	proxy.Infof(sb.String())
	sb.Reset()

	tools.LoadKey(Cfg.Key, &Key)
	proxy.Debugf("劫持私钥 %s 加载成功", Cfg.Key)

	loadCfg(Cfg.Headers, &Headers)
	proxy.Debugf("响应头黑名单 %s 加载成功", Cfg.Headers)

	loadRules(Cfg.QRCode, &QRCode)
	proxy.Debugf("链接劫持规则 %s 加载成功", Cfg.QRCode)

	loadFile(Cfg.Image, &Image)
	proxy.Debugf("劫持图片 %s 加载成功", Cfg.Image)

	loadFile(Cfg.Html, &Html)
	proxy.Debugf("劫持H5代码 %s 加载成功", Cfg.Html)

	loadTargets(Cfg.Targets, &Targets)
	proxy.Debugf("劫持目标地址 %s 加载成功", Cfg.Targets)

	openFile(Cfg.Template, &Template)
	proxy.Debugf("报告模板 %s 加载成功", Cfg.Template)

	proxy.Info("全部配置加载完成")
}

func loadFile(path string, filePtr *[]byte) {
	file, err := os.ReadFile(path)
	if err != nil {
		proxy.Fatal(err)
	}
	*filePtr = file
}

func openFile(path string, filePtr **os.File) {
	file, err := os.Open(path)
	if err != nil {
		proxy.Fatal(err)
	}
	*filePtr = file
}
