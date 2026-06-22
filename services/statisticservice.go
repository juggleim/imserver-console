package services

import (
	"github.com/juggleim/imserver-console/dbs"
)

type StatType int

const (
	StatType_Up       StatType = 1
	StatType_Dispatch StatType = 2
	StatType_Down     StatType = 3

	connectType_Connect       = 0
	connectType_ChrmConnect   = 1
	connectType_ChrmConnect_2 = 2
)

var oneDay int64 = 24 * 60 * 60 * 1000

type Statistics struct {
	Items          []interface{} `json:"items"`
	TotalUserCount *int64        `json:"total_user_count,omitempty"`
}

type StatisticMsgItem struct {
	Count    int64 `json:"count"`
	TimeMark int64 `json:"time_mark"`
}

type MsgStatistics struct {
	MsgUp       *Statistics `json:"msg_up,omitempty"`
	MsgDown     *Statistics `json:"msg_down,omitempty"`
	MsgDispatch *Statistics `json:"msg_dispatch,omitempty"`
}

type UserActivityItem struct {
	Count    int64 `json:"count"`
	TimeMark int64 `json:"time_mark"`
}

type ConcurrentConnectItem struct {
	TimeMark int64 `json:"time_mark"`
	Count    int64 `json:"count"`
}

func QryMsgStatistic(appkey string, statTypes []StatType, channelType int, start, end int64) *MsgStatistics {
	ret := &MsgStatistics{}
	intStateTypes := []int{}
	for _, stateType := range statTypes {
		intStateTypes = append(intStateTypes, int(stateType))
		switch stateType {
		case StatType_Up:
			ret.MsgUp = &Statistics{Items: []interface{}{}}
		case StatType_Down:
			ret.MsgDown = &Statistics{Items: []interface{}{}}
		case StatType_Dispatch:
			ret.MsgDispatch = &Statistics{Items: []interface{}{}}
		}
	}
	dao := dbs.MsgStatDao{}
	list := dao.QryStats(appkey, intStateTypes, channelType, start, end)
	for _, item := range list {
		statItem := &StatisticMsgItem{
			Count:    item.Count,
			TimeMark: item.TimeMark,
		}
		switch item.StatType {
		case int(StatType_Up):
			ret.MsgUp.Items = append(ret.MsgUp.Items, statItem)
		case int(StatType_Down):
			ret.MsgDown.Items = append(ret.MsgDown.Items, statItem)
		case int(StatType_Dispatch):
			ret.MsgDispatch.Items = append(ret.MsgDispatch.Items, statItem)
		}
	}
	return ret
}

func QryUserActivities(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	dao := dbs.DailyActivityDao{}
	list := dao.QryStats(appkey, start, end)
	for _, item := range list {
		ret.Items = append(ret.Items, &UserActivityItem{
			TimeMark: item.TimeMark,
			Count:    item.Count,
		})
	}
	return ret
}

func QryUserRegiste(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	timeMarks := []int64{}
	for s := start / oneDay * oneDay; s <= end; {
		timeMarks = append(timeMarks, s)
		s = s + oneDay
	}
	userDao := dbs.UserDao{}
	for _, timemark := range timeMarks {
		ret.Items = append(ret.Items, &UserActivityItem{
			TimeMark: timemark,
			Count:    userDao.CountByTime(appkey, timemark, timemark+oneDay),
		})
	}
	totalCount := userDao.Count(appkey)
	if totalCount > 0 {
		tc := totalCount
		ret.TotalUserCount = &tc
	}
	return ret
}

func QryConnect(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	dao := dbs.ConnectCountDao{}
	list := dao.QryStats(appkey, connectType_Connect, start, end)
	for _, item := range list {
		ret.Items = append(ret.Items, &ConcurrentConnectItem{
			TimeMark: item.TimeMark,
			Count:    item.Count,
		})
	}
	return ret
}

func QryMaxConnect(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	timeMarks := []int64{}
	for s := start / oneDay * oneDay; s <= end; {
		timeMarks = append(timeMarks, s)
		s = s + oneDay
	}
	dao := dbs.ConnectCountDao{}
	for _, timemark := range timeMarks {
		var count int64
		if item := dao.MaxByTime(appkey, connectType_Connect, timemark, timemark+oneDay); item != nil {
			count = item.Count
		}
		ret.Items = append(ret.Items, &ConcurrentConnectItem{
			TimeMark: timemark,
			Count:    count,
		})
	}
	return ret
}

func QryChrmConnect(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	dao := dbs.ConnectCountDao{}
	list := dao.QryStats(appkey, connectType_ChrmConnect, start, end)
	for _, item := range list {
		ret.Items = append(ret.Items, &ConcurrentConnectItem{
			TimeMark: item.TimeMark,
			Count:    item.Count,
		})
	}
	return ret
}

func QryMaxChrmConnect(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	timeMarks := []int64{}
	for s := start / oneDay * oneDay; s <= end; {
		timeMarks = append(timeMarks, s)
		s = s + oneDay
	}
	dao := dbs.ConnectCountDao{}
	for _, timemark := range timeMarks {
		var count int64
		if item := dao.MaxByTime(appkey, connectType_ChrmConnect, timemark, timemark+oneDay); item != nil {
			count = item.Count
		}
		ret.Items = append(ret.Items, &ConcurrentConnectItem{
			TimeMark: timemark,
			Count:    count,
		})
	}
	return ret
}

func QryMaxChrmConnectV2(appkey string, start, end int64) *Statistics {
	ret := &Statistics{
		Items: []interface{}{},
	}
	timeMarks := []int64{}
	for s := start / oneDay * oneDay; s <= end; {
		timeMarks = append(timeMarks, s)
		s = s + oneDay
	}
	dao := dbs.ConnectCountDao{}
	for _, timemark := range timeMarks {
		var count int64
		if item := dao.MaxByTime(appkey, connectType_ChrmConnect_2, timemark, timemark+oneDay); item != nil {
			count = item.Count
		}
		ret.Items = append(ret.Items, &ConcurrentConnectItem{
			TimeMark: timemark,
			Count:    count,
		})
	}
	return ret
}
