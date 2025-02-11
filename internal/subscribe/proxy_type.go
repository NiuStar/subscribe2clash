package subscribe

const (
	SSRServer = iota
	SSRPort
	SSRProtocol
	SSRCipher
	SSROBFS
	SSRSuffix
)

/*
地址（address）:  portainer.mukj.cn
端口（port）：  443
[::]:443
用户 ID（UUID）： a5b1b0a1-6232-4c76-aea0-ab0480fe2fc6
加密方式（security）： none
传输协议（network）： ws
伪装类型（type）： none
路径（path）： /46019c8aacd9/
底层传输安全： tls
*/
type Vmess struct {
	Add  string `json:"add"`
	Aid  any    `json:"aid"`
	Host string `json:"host"`
	ID   string `json:"id"`
	Net  string `json:"net"`
	Path string `json:"path"`
	Port any    `json:"port"`
	PS   string `json:"ps"`
	TLS  string `json:"tls"`
	Type string `json:"type"`
	V    any    `json:"v"`
}

type ClashVmess struct {
	Name           string    `json:"name,omitempty"`
	Type           string    `json:"type,omitempty"`
	Server         string    `json:"server,omitempty"`
	Port           string    `json:"port,omitempty"`
	UUID           string    `json:"uuid,omitempty"`
	AlterID        string    `json:"alterId,omitempty"`
	Cipher         string    `json:"cipher,omitempty"`
	TLS            bool      `json:"tls,omitempty"`
	Network        string    `json:"network,omitempty"`
	WSOpts         WSOptions `json:"ws-opts,omitempty"`
	SkipCertVerify bool      `json:"skip-cert-verify,omitempty"`
}

type WSOptions struct {
	Path                string            `json:"path,omitempty"`
	Headers             map[string]string `json:"headers,omitempty"`
	MaxEarlyData        int               `json:"max-early-data,omitempty"`
	EarlyDataHeaderName string            `json:"early-data-header-name,omitempty"`
}

type ClashRSSR struct {
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Server        string `json:"server,omitempty"`
	Port          string `json:"port,omitempty"`
	Password      string `json:"password,omitempty"`
	Cipher        string `json:"cipher,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	ProtocolParam string `json:"protocol-param,omitempty"`
	OBFS          string `json:"obfs,omitempty"`
	OBFSParam     string `json:"obfs-param,omitempty"`
}

type ClashSS struct {
	Name       string      `json:"name,omitempty"`
	Type       string      `json:"type,omitempty"`
	Server     string      `json:"server,omitempty"`
	Port       string      `json:"port,omitempty"`
	Password   string      `json:"password,omitempty"`
	Cipher     string      `json:"cipher,omitempty"`
	Plugin     string      `json:"plugin,omitempty"`
	PluginOpts *PluginOpts `json:"plugin-opts,omitempty"`
}

type PluginOpts struct {
	Mode           string `json:"mode"`
	Host           string `json:"host,omitempty"`
	Tls            bool   `json:"tls,omitempty"`
	Path           string `json:"path,omitempty"`
	Mux            bool   `json:"mux,omitempty"`
	SkipCertVerify bool   `json:"skip-cert-verify,omitempty"`
}

type SSD struct {
	Airport      string  `json:"airport"`
	Port         int     `json:"port"`
	Encryption   string  `json:"encryption"`
	Password     string  `json:"password"`
	TrafficUsed  float64 `json:"traffic_used"`
	TrafficTotal float64 `json:"traffic_total"`
	Expiry       string  `json:"expiry"`
	URL          string  `json:"url"`
	Servers      []struct {
		ID            int     `json:"id"`
		Server        string  `json:"server"`
		Ratio         float64 `json:"ratio"`
		Remarks       string  `json:"remarks"`
		Port          string  `json:"port"`
		Encryption    string  `json:"encryption"`
		Password      string  `json:"password"`
		Plugin        string  `json:"plugin"`
		PluginOptions string  `json:"plugin_options"`
	} `json:"servers"`
}

type Clash struct {
	Port      int `yaml:"port"`
	SocksPort int `yaml:"socks-port"`
	// RedirPort          int                      `yaml:"redir-port"`
	// Authentication     []string                 `yaml:"authentication"`
	AllowLan           bool   `yaml:"allow-lan"`
	Mode               string `yaml:"mode"`
	LogLevel           string `yaml:"log-level"`
	ExternalController string `yaml:"external-controller"`
	// ExternalUI         string                   `yaml:"external-ui"`
	// Secret             string                   `yaml:"secret"`
	// Experimental       map[string]interface{} 	`yaml:"experimental"`
	Proxies           []map[string]any `yaml:"proxies"`
	ProxyGroups       []map[string]any `yaml:"proxy-groups"`
	Rule              []string         `yaml:"rules"`
	CFWByPass         []string         `yaml:"cfw-bypass"`
	CFWLatencyTimeout int              `yaml:"cfw-latency-timeout"`

	path       string
	rawProxies []any
}

type Trojan struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Server         string `json:"server"`
	Password       string `json:"password"`
	Sni            string `json:"sni,omitempty"`
	Port           string `json:"port"`
	SkipCertVerify bool   `json:"skip-cert-verify"`
}

type Https struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Server         string `json:"server"`
	UserName       string `json:"username"`
	Password       string `json:"password"`
	Sni            string `json:"sni,omitempty"`
	Tls            bool   `json:"tls,omitempty"`
	Port           string `json:"port"`
	SkipCertVerify bool   `json:"skip-cert-verify"`
}

type Vless struct {
	UUID       string `json:"uuid"`
	Host       string `json:"host"`
	Path       string `json:"path"`
	Encryption string `json:"encryption"`
	Security   string `json:"security"`
	SNI        string `json:"sni"`
	Type       string `json:"type"`
	FP         string `json:"fp"`
	Port       int    `json:"port"`
}
