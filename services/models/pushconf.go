package models

type Platform string
type PushChannel string

const (
	Platform_Android Platform = "Android"
	Platform_IOS     Platform = "iOS"
	Platform_Web     Platform = "Web"
	Platform_PC      Platform = "PC"
	Platform_Harmony Platform = "Harmony"
	Platform_Server  Platform = "Server"
	Platform_Bot     Platform = "Bot"

	PushChannel_Apple  PushChannel = "Apple"
	PushChannel_Huawei PushChannel = "Huawei"
	PushChannel_Xiaomi PushChannel = "Xiaomi"
	PushChannel_OPPO   PushChannel = "Oppo"
	PushChannel_VIVO   PushChannel = "Vivo"
	PushChannel_Jpush  PushChannel = "Jpush"
	PushChannel_FCM    PushChannel = "FCM"
	PushChannel_Meizu  PushChannel = "Meizu"
	PushChannel_HONOR  PushChannel = "Honor"
	PushChannel_Getui  PushChannel = "Getui"
)

type AndroidPushConf struct {
	AppKey      string         `json:"app_key"`
	PushChannel string         `json:"push_channel,omitempty"`
	Package     string         `json:"package"`
	Extra       map[string]any `json:"extra"`
}

type HuaweiPushConf struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func (conf *HuaweiPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppSecret != ""
}

type XiaomiPushConf struct {
	AppSecret string `json:"app_secret"`
	ChannelId string `json:"channel_id"`
}

func (conf *XiaomiPushConf) Valid() bool {
	return conf.AppSecret != ""
}

type OppoPushConf struct {
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
	ChannelId    string `json:"channel_id"`
}

func (conf *OppoPushConf) Valid() bool {
	return conf.AppKey != "" && conf.MasterSecret != ""
}

type VivoPushConf struct {
	AppId     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func (conf *VivoPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppKey != "" && conf.AppSecret != ""
}

type JPushConf struct {
	AppKey       string        `json:"app_key"`
	MasterSecret string        `json:"master_secret"`
	Options      *JPushOptions `json:"options,omitempty"`
}
type JPushOptions struct {
	Classification    int                     `json:"classification,omitempty"`
	ThirdPartyChannel *JPushThirdPartyChannel `json:"third_party_channel,omitempty"`
}

type JPushThirdPartyChannel struct {
	Huawei *JPushHuaweiChannel `json:"huawei,omitempty"`
	Xiaomi *JPushXiaomiChannel `json:"xiaomi,omitempty"`
	Honor  *JPushHonorChannel  `json:"honor,omitempty"`
	Oppo   *JPushOppoChannel   `json:"oppo,omitempty"`
	Vivo   *JPushVivoChannel   `json:"vivo,omitempty"`
	Meizu  *JPushMeizuChannel  `json:"meizu,omitempty"`
}

type JPushHuaweiChannel struct {
	Importance string `json:"importance,omitempty"`
	Category   string `json:"category,omitempty"`
}

type JPushXiaomiChannel struct {
	ChannelId string `json:"channel_id,omitempty"`
}

type JPushHonorChannel struct {
	Importance string `json:"importance,omitempty"`
}

type JPushOppoChannel struct {
	ChannelId   string `json:"channel_id,omitempty"`
	Category    string `json:"category,omitempty"`
	NotifyLevel int    `json:"notify_level,omitempty"`
}

type JPushVivoChannel struct {
	Distribution string `json:"distribution,omitempty"`
	Category     string `json:"category,omitempty"`
}

type JPushMeizuChannel struct {
	Distribution string `json:"distribution,omitempty"`
}

func (conf *JPushConf) Valid() bool {
	return conf.AppKey != "" && conf.MasterSecret != ""
}

type HonorPushConf struct {
	AppId     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func (conf *HonorPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppKey != "" && conf.AppSecret != ""
}

type GetuiPushConf struct {
	AppId        string `json:"app_id"`
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

func (conf *GetuiPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppKey != "" && conf.MasterSecret != ""
}
