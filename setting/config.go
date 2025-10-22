package setting

import (
	"github.com/vpxuser/proxy"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	LogLevel       proxy.Level `yaml:"logLevel"`
	Hotspot        bool        `yaml:"hotspot"`
	Skip           bool        `yaml:"skip"`
	Report         bool        `yaml:"report"`
	Cert           string      `yaml:"cert"`
	Key            string      `yaml:"key"`
	Host           string      `yaml:"host"`
	Targets        string      `yaml:"targets"`
	Headers        string      `yaml:"headers"`
	QRCode         string      `yaml:"qrcode"`
	ImageURL       string      `yaml:"imageURL"`
	QRCodeURL      string      `yaml:"qrcodeURL"`
	Image          string      `yaml:"image"`
	Html           string      `yaml:"html"`
	Template       string      `yaml:"template"`
	CaptureProxy   string      `yaml:"captureProxy"`
	CaptureWifi    string      `yaml:"captureWifi"`
	FileServerHost string      `yaml:"fileServerHost"`
}

func loadCfg(path string, cfg any) {
	file, err := os.ReadFile(path)
	if err != nil {
		proxy.Fatal(err)
	}

	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		proxy.Fatal(err)
	}
}
