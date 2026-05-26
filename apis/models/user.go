package models

type User struct {
	UserId   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Pinyin   string `json:"pinyin"`
	UserType int    `json:"user_type"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Status   int32  `json:"status"`

	CreatedTime int64 `json:"created_time"`
}

type Users struct {
	Items  []*User `json:"items"`
	Offset string  `json:"offset"`
}

type BanUsersReq struct {
	AppKey string     `json:"app_key"`
	Items  []*BanUser `json:"items"`
}

type BanUser struct {
	UserId        string `json:"user_id"`
	EndTime       int64  `json:"end_time"`
	EndTimeOffset int64  `json:"end_time_offset"`
}
