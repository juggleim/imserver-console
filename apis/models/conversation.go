package models

type GlobalConversation struct {
	Sender      *User  `json:"sender"`
	Receiver    *User  `json:"receiver"`
	Group       *Group `json:"group"`
	ChannelType int    `json:"channel_type"`

	Time int64 `json:"time"`
}

type GlobalConversations struct {
	Items []*GlobalConversation `json:"items"`
}
