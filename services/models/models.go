package models

type AddressConf struct {
	Default   []string          `json:"default"`
	NodeConfs map[string]string `json:"confs"`
}

type Apps struct {
	Items   []*SimpleApp `json:"items"`
	HasMore bool         `json:"has_more"`
	Offset  string       `json:"offset"`
}
type SimpleApp struct {
	AppKey       string `json:"app_key"`
	AppType      int    `json:"app_type"`
	AppName      string `json:"app_name"`
	MaxUserCount int    `json:"max_user_count"`
	CurUserCount int64  `json:"cur_user_count"`
	CreatedTime  int64  `json:"created_time"`
	EndedTime    int64  `json:"ended_time"`
}

type Accounts struct {
	Items   []*Account `json:"items"`
	HasMore bool       `json:"has_more"`
	Offset  string     `json:"offset"`
}
type Account struct {
	Account       string `json:"account"`
	State         int    `json:"state"`
	CreatedTime   int64  `json:"created_time"`
	UpdatedTime   int64  `json:"updated_time"`
	ParentAccount string `json:"parent_account"`
	// RoleId        int    `json:"role_id"`
	RoleType int `json:"role_type"`
}

type AppInfo struct {
	AppType     int    `json:"app_type"`
	AppName     string `json:"app_name"`
	CreatedTime int64  `json:"created_time"`
	UpdateTime  int64  `json:"updated_time"`

	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	AppStatus int    `json:"app_status"`
	Alias     string `json:"alias"`

	MaxUserCount int   `json:"max_user_count"`
	CurUserCount int64 `json:"cur_user_count"`

	RestrictedFields *RestrictedFields `json:"restricted_fields"`
	ConfigFields     map[string]string `json:"config_fields"`

	ExpiredTime int64 `json:"expired_time"`

	LicenseConf string `json:"license_conf,omitempty"`
}

type AppInfoResp struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data *AppInfo `json:"data,omitempty"`
}

type RestrictedFields struct {
	MaxUserCount int32 `json:"max_user_count"`
}

type ConfigItem struct {
	Key   string      `json:"key"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type EventSubConfigObj struct {
	EventSubUrl  string `json:"event_sub_url"`
	EventSubAuth string `json:"event_sub_auth"`
}

type EventSubSwitchObj struct {
	PrivateMsgSubSwitch  int `json:"private_msg_sub_switch"`
	GroupMsgSubSwitch    int `json:"group_msg_sub_switch"`
	ChatroomMsgSubSwitch int `json:"chatroom_msg_sub_switch"`
	OnlineSubSwitch      int `json:"online_sub_switch"`
	OfflineSubSwitch     int `json:"offline_sub_switch"`
}

type ZegoConfigObj struct {
	AppId  int64  `json:"app_id"`
	Secret string `json:"secret"`
}

type AgoraConfigObj struct {
	AppId          string `json:"app_id"`
	AppCertificate string `json:"app_certificate"`
}

type LivekitConfigObj struct {
	AppKey     string `json:"app_key"`
	AppSecret  string `json:"app_secret"`
	ServiceUrl string `json:"service_url"`
}

type ActiveAppReq struct {
	License string `json:"license"`
}
