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

type fakeAppFinder struct {
	apps map[string]*dbs.AppInfoDao
}

func (f fakeAppFinder) FindByAppkey(appkey string) *dbs.AppInfoDao {
	return f.apps[appkey]
}

type fakeAppNavStore struct {
	items     map[string]*dbs.AppNavDao
	findError error
	upsertErr error
	ensureErr error
}

func (f *fakeAppNavStore) FindByAppkey(appkey string) (*dbs.AppNavDao, error) {
	if f.findError != nil {
		return nil, f.findError
	}
	item, ok := f.items[appkey]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *item
	return &copy, nil
}

func (f *fakeAppNavStore) FindByAliasNo(aliasNo string) (*dbs.AppNavDao, error) {
	for _, item := range f.items {
		if item.AliasNo == aliasNo {
			copy := *item
			return &copy, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeAppNavStore) UpsertAlias(appkey string, aliasNo string) error {
	if f.upsertErr != nil {
		return f.upsertErr
	}
	f.items[appkey] = &dbs.AppNavDao{AppKey: appkey, AliasNo: aliasNo}
	return nil
}

func (f *fakeAppNavStore) EnsureNextAlias(appkey string) (string, error) {
	if f.ensureErr != nil {
		return "", f.ensureErr
	}
	if item, ok := f.items[appkey]; ok && item.AliasNo != "" && item.AliasNo != "0" {
		return item.AliasNo, nil
	}
	nextAlias := "100000"
	f.items[appkey] = &dbs.AppNavDao{AppKey: appkey, AliasNo: nextAlias}
	return nextAlias, nil
}

func TestQueryAppAlias(t *testing.T) {
	t.Run("defaults to zero when no row exists", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
		if got := queryAppAlias("app-1", store); got != "0" {
			t.Fatalf("got %q, want 0", got)
		}
	})

	t.Run("returns the stored alias", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{
			"app-1": {AppKey: "app-1", AliasNo: "customer-a"},
		}}
		if got := queryAppAlias("app-1", store); got != "customer-a" {
			t.Fatalf("got %q, want customer-a", got)
		}
	})
}

func TestUpdateAppAliasUpserts(t *testing.T) {
	store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
	apps := fakeAppFinder{apps: map[string]*dbs.AppInfoDao{
		"app-1": {AppKey: "app-1"},
	}}

	if code := updateAppAlias("app-1", " customer-a ", apps, store); code != errs.AdminErrorCode_Success {
		t.Fatalf("insert returned code %d", code)
	}
	if got := store.items["app-1"].AliasNo; got != "customer-a" {
		t.Fatalf("inserted alias %q, want customer-a", got)
	}

	if code := updateAppAlias("app-1", "customer-b", apps, store); code != errs.AdminErrorCode_Success {
		t.Fatalf("update returned code %d", code)
	}
	if got := store.items["app-1"].AliasNo; got != "customer-b" {
		t.Fatalf("updated alias %q, want customer-b", got)
	}
}

func TestUpdateAppAliasValidation(t *testing.T) {
	apps := fakeAppFinder{apps: map[string]*dbs.AppInfoDao{
		"app-1": {AppKey: "app-1"},
	}}

	cases := []struct {
		name  string
		alias string
		want  errs.AdminErrorCode
	}{
		{name: "empty", alias: " ", want: errs.AdminErrorCode_ParamError},
		{name: "too long", alias: strings.Repeat("a", MaxAppAliasLength+1), want: errs.AdminErrorCode_ParamError},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
			if got := updateAppAlias("app-1", tc.alias, apps, store); got != tc.want {
				t.Fatalf("got code %d, want %d", got, tc.want)
			}
		})
	}

	store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
	if got := updateAppAlias("missing", "alias", apps, store); got != errs.AdminErrorCode_AppNotExist {
		t.Fatalf("missing app got code %d", got)
	}
}

func TestUpdateAppAliasRejectsDuplicateAndWriteErrors(t *testing.T) {
	logs.SetLogger(logrus.New(), logrus.New())
	apps := fakeAppFinder{apps: map[string]*dbs.AppInfoDao{
		"app-1": {AppKey: "app-1"},
	}}
	store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{
		"app-2": {AppKey: "app-2", AliasNo: "used-alias"},
	}}
	if got := updateAppAlias("app-1", "used-alias", apps, store); got != errs.AdminErrorCode_UpdAppFail {
		t.Fatalf("duplicate alias got code %d", got)
	}

	store = &fakeAppNavStore{
		items:     map[string]*dbs.AppNavDao{},
		upsertErr: errors.New("write failed"),
	}
	if got := updateAppAlias("app-1", "new-alias", apps, store); got != errs.AdminErrorCode_UpdAppFail {
		t.Fatalf("write error got code %d", got)
	}
}

func TestEnsureAppAlias(t *testing.T) {
	logs.SetLogger(logrus.New(), logrus.New())

	t.Run("assigns first automatic alias", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
		code, alias := ensureAppAlias("app-1", store)
		if code != errs.AdminErrorCode_Success || alias != "100000" {
			t.Fatalf("got code=%d alias=%q", code, alias)
		}
	})

	t.Run("reuses existing alias", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{
			"app-1": {AppKey: "app-1", AliasNo: "100008"},
		}}
		code, alias := ensureAppAlias("app-1", store)
		if code != errs.AdminErrorCode_Success || alias != "100008" {
			t.Fatalf("got code=%d alias=%q", code, alias)
		}
	})

	t.Run("rejects empty app key", func(t *testing.T) {
		store := &fakeAppNavStore{items: map[string]*dbs.AppNavDao{}}
		code, alias := ensureAppAlias(" ", store)
		if code != errs.AdminErrorCode_ParamError || alias != "" {
			t.Fatalf("got code=%d alias=%q", code, alias)
		}
	})

	t.Run("reports storage failure", func(t *testing.T) {
		store := &fakeAppNavStore{
			items:     map[string]*dbs.AppNavDao{},
			ensureErr: errors.New("allocate failed"),
		}
		code, alias := ensureAppAlias("app-1", store)
		if code != errs.AdminErrorCode_UpdAppFail || alias != "" {
			t.Fatalf("got code=%d alias=%q", code, alias)
		}
	})
}
