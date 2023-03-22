package boot

import (
	"subscribe2clash/internal/acl"
	"subscribe2clash/internal/global"
)

func Options() []acl.GenOption {
	var options []acl.GenOption
	if global.BaseFile != "" {
		options = append(options, acl.WithBaseFile(global.BaseFile))
	}
	if global.OutputFile != "" {
		options = append(options, acl.WithOutputFile(global.OutputFile))
	}
	return options
}
