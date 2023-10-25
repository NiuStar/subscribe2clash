package subscribe

import (
	"alicode.yjkj.ink/yjkj.ink/work/utils/uuid"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"subscribe2clash/internal/xbase64"
)

const (
	ssrPrefix    = "ssr://"
	vmessPrefix  = "vmess://"
	vlessPrefix  = "vless://"
	ssPrefix     = "ss://"
	trojanPrefix = "trojan://"
	httpsPrefix  = "https://"
	httpPrefix   = "http://"
)

var (
	//ssReg      = regexp.MustCompile(`(?m)ss://(\w+)@([^:]+):(\d+)\?plugin=([^;]+);\w+=(\w+)(?:;obfs-host=)?([^#]+)?#(.+)`)
	ssReg2 = regexp.MustCompile(`(.+):(.+)@(.+):(\d+)`)
	//ssReg2 = regexp.MustCompile(`(?m)^ss://(.+):(.+)@(.+):(\d+)\?(.+)#(.+)$`)

	ssReg = regexp.MustCompile(`(?m)^ss://(.+)\?(.+)#(.+)$`)

	trojanReg = regexp.MustCompile(`(?m)^trojan://(.+)@(.+):(\d+)\?(.+)#(.+)$`)

	httpsReg = regexp.MustCompile(`(.+):(.+)@(.+):(\d+)`)
)

func ParseProxy(contentSlice []string) []any {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
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

	case strings.HasPrefix(proxy, vlessPrefix):
		return vlessConf(proxy)
	case strings.HasPrefix(proxy, ssPrefix):
		return ssConf(proxy)
	case strings.HasPrefix(proxy, trojanPrefix):
		return trojanConf(proxy)
	case strings.HasPrefix(proxy, httpsPrefix):
		return httpsConf(proxy)
	case strings.HasPrefix(proxy, httpPrefix):
		return httpsConf(proxy)
	default:
		fmt.Println(proxy)
	}

	return nil
}
func vlessConf(vlessURL string) *ClashVless {
	u, err := url.Parse(vlessURL)
	if err != nil {
		fmt.Println("URL解析失败：", err)
		return nil
	}

	// 提取vless链接中的参数

	port, _ := strconv.Atoi(u.Port())
	if port == 0 {
		port = 80
	}
	vlessConfig := &Vless{
		UUID:       u.User.Username(),
		Host:       u.Hostname(),
		Path:       u.Query().Get("path"),
		Encryption: u.Query().Get("encryption"),
		Security:   u.Query().Get("security"),
		SNI:        u.Query().Get("sni"),
		Type:       u.Query().Get("type"),
		FP:         u.Query().Get("fp"),
		Port:       port,
	}
	// 解析vless链接中的名称（#后面的部分）
	parts := strings.Split(u.Fragment, "#")
	name := parts[len(parts)-1]
	clashProxyConfig := &ClashVless{
		Name:    name,
		Type:    "vless",
		Server:  vlessConfig.Host,
		Port:    vlessConfig.Port,
		UUID:    vlessConfig.UUID,
		AlterID: 64, // 根据您的设置进行修改
		Cipher:  "none",
		TLS:     true,
		SNI:     vlessConfig.SNI,
		Network: vlessConfig.Type,
		WSPath:  vlessConfig.Path,
	}

	return clashProxyConfig
}

type ClashVless struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Server  string `json:"server"`
	Port    int    `json:"port"`
	UUID    string `json:"uuid"`
	AlterID int    `json:"alterId"`
	Cipher  string `json:"cipher"`
	TLS     bool   `json:"tls"`
	SNI     string `json:"sni"`
	Network string `json:"network"`
	WSPath  string `json:"ws-path"`
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

	rawSSRConfig, err := base64.RawStdEncoding.DecodeString(findStr[1])
	if err != nil {
		return nil
	}

	//s = strings.ReplaceAll(s, findStr[1], string(rawSSRConfig))
	findStr2 := ssReg2.FindStringSubmatch(string(rawSSRConfig))

	ss := &ClashSS{}
	ss.Type = "ss"
	ss.Cipher = findStr2[1]
	ss.Password = findStr2[2]
	ss.Server = findStr2[3]
	ss.Port = findStr2[4]

	ss.Name = findStr[3]

	query := findStr[2]
	queryMap, err := url.ParseQuery(findStr[2])

	if err == nil {
		for k, v := range queryMap {
			ss.Plugin = k
			p := new(PluginOpts)
			switch {
			case strings.Contains(ss.Plugin, "obfs"):
				ss.Plugin = "obfs"
				p.Mode = queryMap["obfs"][0]
				if strings.Contains(query, "obfs-host=") {
					p.Host = queryMap["obfs-host"][0]
				}
			case ss.Plugin == "v2ray-plugin":
				pluginData, _ := base64.RawStdEncoding.DecodeString(v[0])
				json.Unmarshal(pluginData, p)
				p.SkipCertVerify = true
			}
			ss.PluginOpts = p
		}

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

func httpsConf(uriString string) *Https {

	s, err := url.PathUnescape(uriString)
	if err != nil {
		return nil
	}

	lists := strings.Split(s, "#")
	uri, err := url.Parse(lists[0])
	if err != nil {
		return nil
	}

	data, _ := base64.StdEncoding.DecodeString(uri.Host)
	findStr := httpsReg.FindStringSubmatch(string(data))

	if len(findStr) < 5 {
		return nil
	}

	http := &Https{
		UserName: findStr[1],
		Password: findStr[2],
		Type:     "http",
		Server:   findStr[3],
		Port:     findStr[4],
	}
	if len(lists) == 2 {
		http.Name = lists[1]
	} else {
		http.Name = http.Server
	}
	http.Sni = uri.Query().Get("peer")
	if uri.Query().Get("allowInsecure") == "1" {
		http.SkipCertVerify = true

	}
	if uri.Scheme == "https" {
		http.Tls = true
	}
	return http
}
