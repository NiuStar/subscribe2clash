package boot

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"subscribe2clash/constant"
	"subscribe2clash/internal/global"
	"subscribe2clash/internal/req"
)

func init() {
	flag.BoolVar(&global.GenerateConfig, "gc", false, "生成clash配置文件")
	flag.StringVar(&global.BaseFile, "b", "", "clash基础配置文件")
	flag.StringVar(&global.RulesFile, "r", "", "路由配置文件")
	flag.StringVar(&global.OutputFile, "o", "./config/acl.yaml", "clash配置文件名")
	flag.StringVar(&global.Listen, "l", "0.0.0.0:8163", "监听地址")
	flag.StringVar(&req.Proxy, "proxy", "", "http代理")
	flag.IntVar(&global.Tick, "t", 6, "规则更新频率（小时）")
	flag.BoolVar(&global.Version, "version", false, "查看版本信息")
	flag.StringVar(&global.SourceLinks, "link", "", "订阅链接")
	flag.StringVar(&global.SourceFile, "file", "", "订阅文件")
	flag.StringVar(&global.ShareLink, "shareLink", "", "shadowShare分享的订阅地址")

	flag.Parse()
	global.Subscribes = make(map[string]string)
	data, _ := os.ReadFile("./config/config.json")
	json.Unmarshal(data, &global.Subscribes)

	fmt.Println(os.Environ())

	envs := os.Environ()
	for _, e := range envs {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		} else {
			global.Subscribes[parts[0]] = parts[1]
		}
	}
	fmt.Println("shadowShare分享的订阅地址:", global.Subscribes)
	global.Subscribes["mianfeifq"] = "https://raw.fastgit.org/mianfeifq/share/main/data2023036.txt"
}

func initFlag() {
	if global.Version {
		fmt.Printf("subscribe2clash %s %s %s %s\n", constant.Version, runtime.GOOS, runtime.GOARCH, constant.BuildTime)
		os.Exit(0)
	}
}
