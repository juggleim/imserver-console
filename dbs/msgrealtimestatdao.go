package dbs

import (
	"sort"

	"github.com/juggleim/imserver-console/commons/dbcommons"
)

const realtimeBucket30sMs = 30 * 1000

type MsgRealtimeStatDao struct {
	ID          int64  `gorm:"primary_key"`
	StatType    int    `gorm:"stat_type"`
	ChannelType int    `gorm:"channel_type"`
	TimeMark    int64  `gorm:"time_mark"`
	Count       int64  `gorm:"count"`
	AppKey      string `gorm:"app_key"`
}

func (stat MsgRealtimeStatDao) TableName() string {
	return "msgrealtimestats"
}

func (stat MsgRealtimeStatDao) QryStats(appkey string, statTypes []int, channelType int, start, end int64) []*MsgRealtimeStatDao {
	var items []*MsgRealtimeStatDao
	err := dbcommons.GetDb().Where("app_key=? and stat_type in (?) and channel_type=? and time_mark>=? and time_mark<=?", appkey, statTypes, channelType, start, end).
		Order("time_mark asc").
		Limit(30000).
		Find(&items).Error
	if err == nil {
		return items
	}
	return []*MsgRealtimeStatDao{}
}

func (stat MsgRealtimeStatDao) QryStatsWithBucket(appkey string, statTypes []int, channelType int, start, end, bucketMs int64) []*MsgRealtimeStatDao {
	list := stat.QryStats(appkey, statTypes, channelType, start, end)
	if bucketMs <= realtimeBucket30sMs {
		return list
	}
	type bucketKey struct {
		statType int
		timeMark int64
	}
	agg := make(map[bucketKey]int64)
	for _, item := range list {
		key := bucketKey{
			statType: item.StatType,
			timeMark: item.TimeMark / bucketMs * bucketMs,
		}
		agg[key] += item.Count
	}
	result := make([]*MsgRealtimeStatDao, 0, len(agg))
	for key, count := range agg {
		result = append(result, &MsgRealtimeStatDao{
			StatType:    key.statType,
			ChannelType: channelType,
			TimeMark:    key.timeMark,
			Count:       count,
			AppKey:      appkey,
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
