package services

import (
	"os"
	"testing"
	"time"

	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/dbcommons"
)

func TestMain(m *testing.M) {
	if os.Getenv("RUN_DB_TEST") == "1" {
		_ = os.Chdir("..")
		if err := configures.InitConfigures(); err != nil {
			panic(err)
		}
		if err := dbcommons.InitMysql(); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

func TestQryMsgRealtimeAllRangesIntegration(t *testing.T) {
	if os.Getenv("RUN_DB_TEST") != "1" {
		t.Skip("set RUN_DB_TEST=1 to run")
	}
	durations := []struct {
		name string
		ms   int64
	}{
		{"15min", 15 * 60 * 1000},
		{"1h", oneHourMs},
		{"6h", 6 * oneHourMs},
		{"12h", 12 * oneHourMs},
		{"1d", oneDay},
		{"3d", threeDays},
	}
	end := time.Now().UnixMilli()
	for _, tc := range durations {
		t.Run(tc.name, func(t *testing.T) {
			start := end - tc.ms
			ret := QryMsgRealtimeStatistic("appkey", []StatType{StatType_Up, StatType_Down, StatType_Dispatch}, 1, start, end)
			if ret.MsgUp == nil || len(ret.MsgUp.Items) == 0 {
				t.Fatalf("expected msg_up items for %s, got %+v", tc.name, ret.MsgUp)
			}
			t.Logf("%s: msg_up=%d bucketMs=%d", tc.name, len(ret.MsgUp.Items), realtimeBucketMs(end-start))
		})
	}
}
