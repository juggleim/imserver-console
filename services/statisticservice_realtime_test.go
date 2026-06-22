package services

import (
	"testing"

	"github.com/juggleim/imserver-console/dbs"
)

func TestRealtimeBucketMs(t *testing.T) {
	cases := []struct {
		duration int64
		want     int64
	}{
		{15 * 60 * 1000, realtimeBucket30s},
		{oneHourMs, realtimeBucket30s},
		{oneHourMs + 1, realtimeBucket1m},
		{6 * oneHourMs, realtimeBucket1m},
		{6*oneHourMs + 1, realtimeBucket5m},
		{12 * oneHourMs, realtimeBucket5m},
		{oneDay, realtimeBucket10m},
		{oneDay + 1, realtimeBucket15m},
	}
	for _, tc := range cases {
		if got := realtimeBucketMs(tc.duration); got != tc.want {
			t.Fatalf("duration=%d got=%d want=%d", tc.duration, got, tc.want)
		}
	}
}

func TestRealtimeAvgPerSecond(t *testing.T) {
	cases := []struct {
		count    int64
		bucketMs int64
		want     float64
	}{
		{60, realtimeBucket30s, 2},
		{120, realtimeBucket1m, 2},
		{300, realtimeBucket5m, 1},
		{10, realtimeBucket30s, 0.33},
	}
	for _, tc := range cases {
		if got := realtimeAvgPerSecond(tc.count, tc.bucketMs); got != tc.want {
			t.Fatalf("count=%d bucketMs=%d got=%v want=%v", tc.count, tc.bucketMs, got, tc.want)
		}
	}
}

func TestAggregateRealtimeStatsOneMinute(t *testing.T) {
	bucket := int64(realtimeBucket1m)
	base := int64(1_782_112_500_000)
	list := []*dbs.MsgRealtimeStatDao{
		{StatType: 1, TimeMark: base, Count: 10},
		{StatType: 1, TimeMark: base + 30000, Count: 20},
		{StatType: 1, TimeMark: base + 60000, Count: 5},
	}
	got := aggregateRealtimeStats(list, bucket)
	if len(got) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(got))
	}
	if got[0].Count != 30 || got[1].Count != 5 {
		t.Fatalf("unexpected counts: %+v", got)
	}
}
