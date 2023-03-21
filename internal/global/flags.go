package global

var (
	GenerateConfig bool
	BaseFile       string
	RulesFile      string
	OutputFile     string
	Listen         string
	Tick           int
	Version        bool

	SourceLinks string
	SourceFile  string
	ShareLink   string //shadowShare分享的订阅

	Subscribes map[string]string
)
