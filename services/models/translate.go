package models

type TransEngineConf struct {
	Channel          string            `json:"channel,omitempty"`
	BdTransEngine    *BdTransEngine    `json:"baidu,omitempty"`
	DeeplTransEngine *DeeplTransEngine `json:"deepl,omitempty"`
}

type BdTransEngine struct {
	AppKey      string `json:"-"`
	ApiKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	accessToken string `json:"-"`
}

type DeeplTransEngine struct {
	AppKey  string `json:"-"`
	AuthKey string `json:"auth_key"`
}
