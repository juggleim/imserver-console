package models

type MailEngineConf struct {
	Channel string                 `json:"channel,omitempty"`
	Ali     map[string]interface{} `json:"ali,omitempty"`
	Engagelab map[string]interface{} `json:"engagelab,omitempty"`
	Neteasy map[string]interface{} `json:"neteasy,omitempty"`
}
