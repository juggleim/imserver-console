package services

import (
	"errors"
	"strings"
	"testing"

	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func TestQueryAppNavUrl(t *testing.T) {
	t.Run("returns empty when no row exists", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
		if got := queryAppNavUrl("app-1", AppNavWsUrlColumn, store); got != "" {
			t.Fatalf("got %q, want empty", got)
		}
	})

	t.Run("returns the stored url per column", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{
			"app-1": {AppKey: "app-1", AliasNo: "100001", WsUrl: "ws://ws.example.com", ApiUrl: "https://api.example.com", AppUrl: "https://app.example.com"},
		}}
		if got := queryAppNavUrl("app-1", AppNavWsUrlColumn, store); got != "ws://ws.example.com" {
			t.Fatalf("ws_url got %q", got)
		}
		if got := queryAppNavUrl("app-1", AppNavApiUrlColumn, store); got != "https://api.example.com" {
			t.Fatalf("api_url got %q", got)
		}
		if got := queryAppNavUrl("app-1", AppNavAppUrlColumn, store); got != "https://app.example.com" {
			t.Fatalf("app_url got %q", got)
		}
	})

	t.Run("returns empty on unknown column", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{
			"app-1": {AppKey: "app-1", WsUrl: "ws://ws.example.com"},
		}}
		if got := queryAppNavUrl("app-1", "admin_url", store); got != "" {
			t.Fatalf("got %q, want empty", got)
		}
	})
}

func TestUpdateAppUrlUpserts(t *testing.T) {
	logs.SetLogger(logrus.New(), logrus.New())
	apps := fakeAppFinder{apps: map[string]*dbs.AppInfoDao{
		"app-1": {AppKey: "app-1"},
	}}

	columns := []struct {
		name   string
		column string
	}{
		{name: "ws", column: AppNavWsUrlColumn},
		{name: "api", column: AppNavApiUrlColumn},
		{name: "app", column: AppNavAppUrlColumn},
	}
	for _, tc := range columns {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
			if code := updateAppUrl("app-1", " https://one.example.com ", tc.column, apps, store); code != errs.AdminErrorCode_Success {
				t.Fatalf("insert returned code %d", code)
			}
			saved := store.items["app-1"]
			switch tc.column {
			case AppNavWsUrlColumn:
				if saved.WsUrl != "https://one.example.com" {
					t.Fatalf("inserted ws_url %q", saved.WsUrl)
				}
			case AppNavApiUrlColumn:
				if saved.ApiUrl != "https://one.example.com" {
					t.Fatalf("inserted api_url %q", saved.ApiUrl)
				}
			case AppNavAppUrlColumn:
				if saved.AppUrl != "https://one.example.com" {
					t.Fatalf("inserted app_url %q", saved.AppUrl)
				}
			}

			if code := updateAppUrl("app-1", "https://two.example.com", tc.column, apps, store); code != errs.AdminErrorCode_Success {
				t.Fatalf("update returned code %d", code)
			}
			saved = store.items["app-1"]
			switch tc.column {
			case AppNavWsUrlColumn:
				if saved.WsUrl != "https://two.example.com" {
					t.Fatalf("updated ws_url %q", saved.WsUrl)
				}
			case AppNavApiUrlColumn:
				if saved.ApiUrl != "https://two.example.com" {
					t.Fatalf("updated api_url %q", saved.ApiUrl)
				}
			case AppNavAppUrlColumn:
				if saved.AppUrl != "https://two.example.com" {
					t.Fatalf("updated app_url %q", saved.AppUrl)
				}
			}
		})
	}
}

func TestUpdateAppUrlValidation(t *testing.T) {
	logs.SetLogger(logrus.New(), logrus.New())
	apps := fakeAppFinder{apps: map[string]*dbs.AppInfoDao{
		"app-1": {AppKey: "app-1"},
	}}

	cases := []struct {
		name string
		url  string
		want errs.AdminErrorCode
	}{
		{name: "empty", url: " ", want: errs.AdminErrorCode_ParamError},
		{name: "too long", url: strings.Repeat("a", MaxAppUrlLength+1), want: errs.AdminErrorCode_ParamError},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
			if got := updateAppUrl("app-1", tc.url, AppNavWsUrlColumn, apps, store); got != tc.want {
				t.Fatalf("got code %d, want %d", got, tc.want)
			}
		})
	}

	t.Run("missing app", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
		if got := updateAppUrl("missing", "https://api.example.com", AppNavWsUrlColumn, apps, store); got != errs.AdminErrorCode_AppNotExist {
			t.Fatalf("missing app got code %d", got)
		}
	})

	t.Run("unknown column", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
		if got := updateAppUrl("app-1", "https://api.example.com", "admin_url", apps, store); got != errs.AdminErrorCode_ParamError {
			t.Fatalf("unknown column got code %d", got)
		}
	})
}

func TestUpdateAppUrlWriteErrors(t *testing.T) {
	logs.SetLogger(logrus.New(), logrus.New())
	apps := fakeAppFinder{apps: map[string]*dbs.AppInfoDao{
		"app-1": {AppKey: "app-1"},
	}}

	store := &fakeAppNavStore{
		items:        map[string]*dbs.AppNavDao{},
		upsertUrlErr: errors.New("write failed"),
	}
	if got := updateAppUrl("app-1", "https://api.example.com", AppNavWsUrlColumn, apps, store); got != errs.AdminErrorCode_UpdAppFail {
		t.Fatalf("write error got code %d", got)
	}

	store = &fakeAppNavStore{
		items:     map[string]*dbs.AppNavDao{},
		findError: gorm.ErrRecordNotFound,
	}
	if got := updateAppUrl("app-1", "https://api.example.com", AppNavWsUrlColumn, apps, store); got != errs.AdminErrorCode_UpdAppFail {
		t.Fatalf("verify error got code %d", got)
	}
}
