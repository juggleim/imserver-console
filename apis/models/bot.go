package models

type Bot struct {
	BotId       string       `json:"bot_id"`
	Nickname    string       `json:"nickname"`
	Avatar      string       `json:"avatar"`
	Pinyin      string       `json:"pinyin"`
	UserType    int          `json:"user_type"`
	BotConf     *BotConf     `json:"bot_conf"`
	BotSettings *BotSettings `json:"bot_settings"`

	CreatedTime int64 `json:"created_time"`
}

type Bots struct {
	Items  []*Bot `json:"items"`
	Offset string `json:"offset"`
}

type BotConf struct {
	BotId    string `json:"bot_id"`
	Url      string `json:"url"`
	ApiKey   string `json:"api_key"`
	IsStream bool   `json:"is_stream"`
}

type BotSettings struct {
	OnlyMentioned bool `json:"only_mentioned"`
}

type BotReq struct {
	AppKey      string       `json:"app_key"`
	BotId       string       `json:"bot_id"`
	Nickname    string       `json:"nickname"`
	Avatar      string       `json:"avatar"`
	Pinyin      string       `json:"pinyin"`
	BotConf     *BotConf     `json:"bot_conf"`
	BotSettings *BotSettings `json:"bot_settings"`
	ExtFields   map[string]string `json:"ext_fields"`
}
