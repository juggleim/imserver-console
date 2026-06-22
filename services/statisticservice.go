package services

import (
	"math"
	"sort"

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
var threeDays int64 = 3 * 24 * 60 * 60 * 1000
var oneHourMs int64 = 60 * 60 * 1000
var sixHoursMs int64 = 6 * oneHourMs
var twelveHoursMs int64 = 12 * oneHourMs

const (
	realtimeBucket30s = 30 * 1000
	realtimeBucket1m  = 60 * 1000
	realtimeBucket5m  = 5 * 60 * 1000
	realtimeBucket10m = 10 * 60 * 1000
	realtimeBucket15m = 15 * 60 * 1000
)

func ThreeDaysMs() int64 {
	return threeDays
}

type Statistics struct {
	Items          []interface{} `json:"items"`
	TotalUserCount *int64        `json:"total_user_count,omitempty"`
}

type StatisticMsgItem struct {
	Count    float64 `json:"count"`
	TimeMark int64   `json:"time_mark"`
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
			Count:    float64(item.Count),
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

func realtimeBucketMs(duration int64) int64 {
	if duration <= 0 {
		return realtimeBucket30s
	}
	if duration <= oneHourMs {
		return realtimeBucket30s
	}
	if duration <= sixHoursMs {
		return realtimeBucket1m
	}
	if duration <= twelveHoursMs {
		return realtimeBucket5m
	}
	if duration <= oneDay {
		return realtimeBucket10m
	}
	return realtimeBucket15m
}

func normalizeEpochMs(ts int64) int64 {
	if ts > 0 && ts < 1_000_000_000_000 {
		return ts * 1000
	}
	return ts
}

func realtimeAvgPerSecond(count int64, bucketMs int64) float64 {
	if bucketMs <= 0 {
		bucketMs = realtimeBucket30s
	}
	avg := float64(count) / (float64(bucketMs) / 1000)
	return math.Round(avg*100) / 100
}

type realtimeAggKey struct {
	statType int
	timeMark int64
}

func aggregateRealtimeStats(list []*dbs.MsgRealtimeStatDao, bucketMs int64) []*dbs.MsgRealtimeStatDao {
	if bucketMs <= realtimeBucket30s {
		return list
	}
	agg := make(map[realtimeAggKey]int64)
	for _, item := range list {
		key := realtimeAggKey{
			statType: item.StatType,
			timeMark: item.TimeMark / bucketMs * bucketMs,
		}
		agg[key] += item.Count
	}
	result := make([]*dbs.MsgRealtimeStatDao, 0, len(agg))
	for key, count := range agg {
		result = append(result, &dbs.MsgRealtimeStatDao{
			StatType: key.statType,
			TimeMark: key.timeMark,
			Count:    count,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].TimeMark != result[j].TimeMark {
			return result[i].TimeMark < result[j].TimeMark
		}
		return result[i].StatType < result[j].StatType
	})
	return result
}

func QryMsgRealtimeStatistic(appkey string, statTypes []StatType, channelType int, start, end int64) *MsgStatistics {
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
	dao := dbs.MsgRealtimeStatDao{}
	start = normalizeEpochMs(start)
	end = normalizeEpochMs(end)
	bucketMs := realtimeBucketMs(end - start)
	list := dao.QryStatsWithBucket(appkey, intStateTypes, channelType, start, end, bucketMs)
	for _, item := range list {
		statItem := &StatisticMsgItem{
			Count:    realtimeAvgPerSecond(item.Count, bucketMs),
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
