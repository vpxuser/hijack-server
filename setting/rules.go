package setting

import (
	"regexp"
)

func loadRules(path string, rules *map[string]*regexp.Regexp) {
	ruleTexts := make(map[string]string)
	loadCfg(path, &ruleTexts)
	for detail, ruleText := range ruleTexts {
		(*rules)[detail] = regexp.MustCompile(ruleText)
	}
}
