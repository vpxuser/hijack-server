package tools

import (
	"github.com/vpxuser/proxy"
	"regexp"
)

var ImageLinkRule = new(regexp.Regexp)

func init() {
	rule, err := regexp.Compile(`(?i)https?:\\?/\\?/[^\s"'<>]*\.(jpg|jpeg|png|gif|webp|bmp|svg)[^\s"'<>]*`)
	if err != nil {
		proxy.Fatal(err)
	}
	ImageLinkRule = rule
}
