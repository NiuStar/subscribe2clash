package boot

import (
	"log"
	"os"
	"sync"
	"time"

	"subscribe2clash/internal/acl"
	"subscribe2clash/internal/clash"
	"subscribe2clash/internal/global"
)

func generateConfig() {
	// 配置文件相关设置
	options := Options()
	g := acl.New(options...)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		ticker := time.NewTicker(time.Duration(global.Tick) * time.Hour)
		for i := 0; true; <-ticker.C {
			g.GenerateConfig()
			if i == 0 {
				wg.Done()
				i++
			}
		}
	}()

	wg.Wait()

	var config string
	var err error
	switch {
	case global.SourceFile != "":
		config, err = clash.Config(clash.File, global.SourceFile)
	case global.SourceLinks != "":
		config, err = clash.Config(clash.Url, global.SourceLinks)
	default:
		return
	}

	if err != nil {
		log.Fatal("生成配置文件内容失败", err)
	}

	// 写入配置文件
	err = os.WriteFile(global.OutputFile, []byte(config), 0644)
	if err != nil {
		log.Fatal("写入配置文件失败", err)
	}

	os.Exit(0)
}
