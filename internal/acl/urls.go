package acl

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/ini.v1"

	"subscribe2clash/internal/global"
)

//go:embed config/default_rules.ini
var defaultRulesConfig []byte

type Rules struct {
	url   string
	rule  string
	local string
}

func GetRules() []Rules {
	var (
		rs  []Rules
		cfg *ini.File
		err error
	)

	if global.RulesFile != "" {
		cfg, err = ini.Load(global.RulesFile)
	} else {
		cfg, err = ini.Load(defaultRulesConfig)
	}
	if err != nil {
		log.Fatal(err)
	}

	host := cfg.Section("").Key("host").String()
	local := cfg.Section("").Key("local").String()

	for _, cfgK := range cfg.Section("rules").Keys() {
		rs = append(rs, Rules{
			url:   fmt.Sprintf("%s/%s", host, cfgK.Name()),
			local: fmt.Sprintf("%s/%s", local, cfgK.Name()),
			rule:  cfgK.Value(),
		})
	}

	return rs
}
