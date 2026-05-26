package models

type Application struct {
	AppId       string `json:"app_id"`
	AppName     string `json:"app_name"`
	AppIcon     string `json:"app_icon"`
	AppDesc     string `json:"app_desc"`
	AppUrl      string `json:"app_url"`
	AppOrder    int    `json:"app_order"`
	CreatedTime int64  `json:"created_time"`
	UpdatedTime int64  `json:"updated_time"`

	AppKey string `json:"app_key,omitempty"`
}

type Applications struct {
	Items  []*Application `json:"items"`
	Offset string         `json:"offset,omitempty"`
	Page   int            `json:"page"`
	Size   int            `json:"size"`
}

type ApplicationIds struct {
	AppKey string   `json:"app_key"`
	AppIds []string `json:"app_ids"`
}
