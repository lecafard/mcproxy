package schemas

type ConfigServer struct {
	Proxy       string `json:"proxy"`
	Description string `json:"description"`
	Startup     string `json:"startup"`
}

type Config struct {
	Listen  string                   `json:"listen"`
	Servers map[string]*ConfigServer `json:"servers"`
	Timeout int                      `json:"timeout"`
}
