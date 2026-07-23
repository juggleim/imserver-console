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
	AppKey          string         `json:"app_key"`
	PushChannel     string         `json:"push_channel,omitempty"`
	Package         string         `json:"package"`
	OriginalPackage string         `json:"original_package,omitempty"`
	Extra           map[string]any `json:"extra"`
}

const PushSecretMask = "********"

type PushConfListItem struct {
	AppKey      string         `json:"app_key"`
	PushChannel string         `json:"push_channel"`
	Package     string         `json:"package"`
	ConfPath    string         `json:"conf_path,omitempty"`
	Extra       map[string]any `json:"extra,omitempty"`
}

type IosPushConfListItem struct {
	AppKey       string `json:"app_key"`
	Package      string `json:"package"`
	IsProduct    int    `json:"is_product"`
	CertPath     string `json:"cert_path,omitempty"`
	CertPwd      string `json:"cert_pwd,omitempty"`
	VoipCertPwd  string `json:"voip_cert_pwd,omitempty"`
	VoipCertPath string `json:"voip_cert_path,omitempty"`
}

type HuaweiPushConf struct {
	AppId      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
	BadgeClass string `json:"badge_class,omitempty"`
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
	BadgeClass   string        `json:"badge_class,omitempty"`
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
	ChannelId       string `json:"channel_id,omitempty"`
	MiTemplateId    string `json:"mi_template_id"`
	MiTemplateParam string `json:"mi_template_param"`
}

type JPushHonorChannel struct {
	Importance string `json:"importance,omitempty"`
}

type JPushOppoChannel struct {
	Distribution             string            `json:"distribution,omitempty"`
	ChannelId                string            `json:"channel_id,omitempty"`
	Category                 string            `json:"category,omitempty"`
	NotifyLevel              int               `json:"notify_level,omitempty"`
	BadgeOperationType       *int              `json:"badge_operation_type,omitempty"`
	PrivateMsgTemplateId     string            `json:"private_msg_template_id,omitempty"`
	PrivateContentParameters map[string]string `json:"private_content_parameters,omitempty"`
	PrivateTitleParameters   map[string]string `json:"private_title_parameters,omitempty"`
}

type JPushVivoChannel struct {
	Distribution string `json:"distribution,omitempty"`
	Category     string `json:"category,omitempty"`
	AddBadge     bool   `json:"add_badge,omitempty"`
}

type JPushMeizuChannel struct {
	Distribution string `json:"distribution,omitempty"`
}

func (conf *JPushConf) Valid() bool {
	return conf.AppKey != "" && conf.MasterSecret != ""
}

type HonorPushConf struct {
	AppId      string `json:"app_id"`
	AppKey     string `json:"app_key"`
	AppSecret  string `json:"app_secret"`
	BadgeClass string `json:"badge_class,omitempty"`
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
