package subscribe

import (
	"alicode.yjkj.ink/yjkj.ink/work/utils/uuid"
	"bufio"
	"encoding/json"
	"fmt"
	"subscribe2clash/internal/xbase64"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	ssrPrefix    = "ssr://"
	vmessPrefix  = "vmess://"
	ssPrefix     = "ss://"
	trojanPrefix = "trojan://"
)

var (
	//ssReg      = regexp.MustCompile(`(?m)ss://(\w+)@([^:]+):(\d+)\?plugin=([^;]+);\w+=(\w+)(?:;obfs-host=)?([^#]+)?#(.+)`)
	ssReg2 = regexp.MustCompile(`(?m)ss://([\-0-9a-z]+):(.+)@(.+):(\d+)(.+)?#(.+)`)
	ssReg  = regexp.MustCompile(`(?m)ss://([/+=\w]+)(@.+)?#(.+)`)

	trojanReg = regexp.MustCompile(`(?m)^trojan://(.+)@(.+):(\d+)\?(.+)#(.+)$`)
)

func ParseProxy(contentSlice []string) []any {
	var proxies []any
	for _, v := range contentSlice {
		// ssd
		if strings.Contains(v, "airport") {
			ssSlice := ssdConf(v)
			for _, ss := range ssSlice {
				if ss.Name != "" {
					proxies = append(proxies, ss)
				}
			}
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(v))
		for scanner.Scan() {
			proxy := parseProxy(scanner.Text())
			if proxy != nil {
				proxies = append(proxies, proxy)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("parse proxy failed, err: %v", err)
		}
	}

	return proxies
}

func subProtocolBody(proxy string, prefix string) string {
	return strings.TrimSpace(proxy[len(prefix):])
	//return proxy
}

func parseProxy(proxy string) any {
	switch {
	case strings.HasPrefix(proxy, ssrPrefix):
		return ssrConf(subProtocolBody(proxy, ssrPrefix))
	case strings.HasPrefix(proxy, vmessPrefix):
		return v2rConf(subProtocolBody(proxy, vmessPrefix))
	case strings.HasPrefix(proxy, ssPrefix):
		return ssConf(proxy)
	case strings.HasPrefix(proxy, trojanPrefix):
		return trojanConf(proxy)
	}

	return nil
}

func v2rConf(s string) *ClashVmess {
	vmconfig, err := xbase64.Base64DecodeStripped(s)
	if err != nil {
		return nil
	}
	vmess := Vmess{}
	err = json.Unmarshal(vmconfig, &vmess)
	if err != nil {
		log.Printf("v2ray config json unmarshal failed, err: %v", err)
		return nil
	}
	if vmess.PS == "澳大利亚 Cranbourne Secondary College" {
		fmt.Println("澳大利亚 Cranbourne Secondary College")
	}
	clashVmess := &ClashVmess{}
	clashVmess.Name = vmess.PS

	clashVmess.Type = "vmess"
	clashVmess.Server = vmess.Add
	switch vmess.Port.(type) {
	case string:
		clashVmess.Port = vmess.Port.(string)
	case int:
		clashVmess.Port = strconv.Itoa(vmess.Port.(int))
	case float64:
		clashVmess.Port = strconv.Itoa(int(vmess.Port.(float64)))
	default:

	}
	_, errUid := uuid.FromString(vmess.ID)
	if errUid == nil {
		clashVmess.UUID = vmess.ID
	} else {
		uid, _ := uuid.NewV4()
		clashVmess.UUID = uid.String()
	}

	switch vmess.Aid.(type) {
	case string:
		clashVmess.AlterID = vmess.Aid.(string)
	case int:
		clashVmess.AlterID = strconv.Itoa(vmess.Aid.(int))
	case float64:
		clashVmess.AlterID = strconv.Itoa(int(vmess.Aid.(float64)))
	default:

	}
	if vmess.Type != "none" && vmess.Type != "" && vmess.Type != "http" {
		clashVmess.Cipher = vmess.Type

	} else {
		clashVmess.Cipher = "auto"

	}

	if strings.EqualFold(vmess.TLS, "tls") {
		clashVmess.TLS = true
	} else {
		clashVmess.TLS = false
	}
	if vmess.Net == "ws" {
		clashVmess.Network = vmess.Net
		clashVmess.WSOpts.Path = vmess.Path
		if vmess.Host != "" {
			clashVmess.WSOpts.Headers = make(map[string]string)
			clashVmess.WSOpts.Headers["Host"] = vmess.Host

		}
	}

	return clashVmess
}

func ssdConf(ssdJson string) []*ClashSS {
	var ssd SSD
	err := json.Unmarshal([]byte(ssdJson), &ssd)
	if err != nil {
		log.Println("ssd json unmarshal err:", err)
		return nil
	}

	var clashSSSlice []*ClashSS
	for _, server := range ssd.Servers {
		options, err := url.ParseQuery(server.PluginOptions)
		if err != nil {
			continue
		}

		var ss = &ClashSS{}
		ss.Type = "ss"
		ss.Name = server.Remarks
		ss.Cipher = server.Encryption
		ss.Password = server.Password
		ss.Server = server.Server
		ss.Port = server.Port
		ss.Plugin = server.Plugin
		ss.PluginOpts = &PluginOpts{
			Mode: options["obfs"][0],
			Host: options["obfs-host"][0],
		}

		switch {
		case strings.Contains(ss.Plugin, "obfs"):
			ss.Plugin = "obfs"
		}

		clashSSSlice = append(clashSSSlice, ss)
	}

	return clashSSSlice
}

func ssrConf(s string) *ClashRSSR {
	rawSSRConfig, err := xbase64.Base64DecodeStripped(s)
	if err != nil {
		return nil
	}
	params := strings.Split(string(rawSSRConfig), `:`)
	if len(params) != 6 {
		return nil
	}
	ssr := &ClashRSSR{}
	ssr.Type = "ssr"
	ssr.Server = params[SSRServer]
	ssr.Port = params[SSRPort]
	ssr.Protocol = params[SSRProtocol]
	ssr.Cipher = params[SSRCipher]
	ssr.OBFS = params[SSROBFS]

	// 如果兼容ss协议，就转换为clash的ss配置
	// https://github.com/Dreamacro/clash
	if ssr.Protocol == "origin" && ssr.OBFS == "plain" {
		switch ssr.Cipher {
		case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm",
			"aes-128-cfb", "aes-192-cfb", "aes-256-cfb",
			"aes-128-ctr", "aes-192-ctr", "aes-256-ctr",
			"rc4-md5", "chacha20", "chacha20-ietf", "xchacha20",
			"chacha20-ietf-poly1305", "xchacha20-ietf-poly1305":
			ssr.Type = "ss"
		}
	}
	suffix := strings.Split(params[SSRSuffix], "/?")
	if len(suffix) != 2 {
		return nil
	}
	passwordBase64 := suffix[0]
	password, err := xbase64.Base64DecodeStripped(passwordBase64)
	if err != nil {
		return nil
	}
	ssr.Password = string(password)

	m, err := url.ParseQuery(suffix[1])
	if err != nil {
		return nil
	}

	for k, v := range m {
		de, err := xbase64.Base64DecodeStripped(v[0])
		if err != nil {
			return nil
		}
		switch k {
		case "obfsparam":
			ssr.OBFSParam = string(de)
			continue
		case "protoparam":
			ssr.ProtocolParam = string(de)
			continue
		case "remarks":
			ssr.Name = string(de)
			continue
		case "group":
			continue
		}
	}

	return ssr
}

func ssConf(s string) *ClashSS {
	s, err := url.PathUnescape(s)
	if err != nil {
		return nil
	}

	findStr := ssReg.FindStringSubmatch(s)
	if len(findStr) < 4 {
		return nil
	}

	rawSSRConfig, err := xbase64.Base64DecodeStripped(findStr[1])
	if err != nil {
		return nil
	}

	s = strings.ReplaceAll(s, findStr[1], string(rawSSRConfig))
	findStr = ssReg2.FindStringSubmatch(s)

	ss := &ClashSS{}
	ss.Type = "ss"
	ss.Cipher = findStr[1]
	ss.Password = findStr[2]
	ss.Server = findStr[3]
	ss.Port = findStr[4]
	ss.Name = findStr[6]

	if findStr[5] != "" && strings.Contains(findStr[5], "plugin") {
		query := findStr[5][strings.Index(findStr[5], "?")+1:]
		queryMap, err := url.ParseQuery(query)
		if err != nil {
			return nil
		}

		ss.Plugin = queryMap["plugin"][0]
		p := new(PluginOpts)
		switch {
		case strings.Contains(ss.Plugin, "obfs"):
			ss.Plugin = "obfs"
			p.Mode = queryMap["obfs"][0]
			if strings.Contains(query, "obfs-host=") {
				p.Host = queryMap["obfs-host"][0]
			}
		case ss.Plugin == "v2ray-plugin":
			p.Mode = queryMap["mode"][0]
			if strings.Contains(query, "host=") {
				p.Host = queryMap["host"][0]
			}
			if strings.Contains(query, "path=") {
				p.Path = queryMap["path"][0]
			}
			p.Mux = strings.Contains(query, "mux")
			p.Tls = strings.Contains(query, "tls")
			p.SkipCertVerify = true
		}
		ss.PluginOpts = p
	}

	return ss
}

func trojanConf(s string) *Trojan {

	s, err := url.PathUnescape(s)
	if err != nil {
		return nil
	}
	if strings.Contains(s, "上海市") {
		fmt.Println()
	}
	findStr := trojanReg.FindStringSubmatch(s)

	if len(findStr) == 6 {
		trojan := &Trojan{
			Name:     findStr[5],
			Type:     "trojan",
			Server:   findStr[2],
			Password: findStr[1],
			Port:     findStr[3],
		}
		values, _ := url.ParseQuery(findStr[4])
		if values != nil {
			if values.Get("sni") != "" {
				trojan.Sni = values.Get("sni")
			}
			if values.Get("allowInsecure") == "1" {
				trojan.SkipCertVerify = true
			} else {
				trojan.SkipCertVerify = false
			}

		}

		return trojan
	}

	return nil
}
