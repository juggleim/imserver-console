package models

type SmsEngineConf struct {
	Channel     string       `json:"channel,omitempty"`
	BdSmsEngine *BdSmsEngine `json:"baidu,omitempty"`
}

type BdSmsEngine struct {
	ApiKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	Endpoint    string `json:"endpoint"`
	Template    string `json:"template"`
	SignatureId string `json:"signature_id"`
}
