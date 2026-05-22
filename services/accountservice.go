package services

import (
	"context"
	"time"

	"github.com/juggleim/imserver-console/commons/caches"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

type RoleType int
type AccountState int

const (
	RoleType_SuperAdmin  RoleType = 0
	RoleType_NormalAdmin RoleType = 1

	AccountState_Normal AccountState = 0
)

type AccountInfo struct {
	Account  string
	RoleType RoleType
	State    AccountState
}

var accountInfoCache *caches.LruCache
var accountLocks *tools.SegmentatedLocks

func init() {
	accountLocks = tools.NewSegmentatedLocks(16)
	accountInfoCache = caches.NewLruCacheWithAddReadTimeout("account_cache", 100, nil, 10*time.Minute, 10*time.Minute)
}

var notExistAccount = &AccountInfo{}

func GetAccountInfo(account string) (*AccountInfo, bool) {
	if val, exist := accountInfoCache.Get(account); exist {
		info := val.(*AccountInfo)
		if info == notExistAccount {
			return nil, false
		}
		return info, true
	} else {
		l := accountLocks.GetLocks(account)
		l.Lock()
		defer l.Unlock()
		if val, exist := accountInfoCache.Get(account); exist {
			info := val.(*AccountInfo)
			if info == notExistAccount {
				return nil, false
			}
			return info, true
		} else {
			dao := dbs.AccountDao{}
			acc, err := dao.FindByAccount(account)
			if err != nil || acc == nil {
				accountInfoCache.Add(account, notExistAccount)
				return nil, false
			}
			info := &AccountInfo{
				Account:  account,
				State:    AccountState(acc.State),
				RoleType: RoleType(acc.RoleType),
			}
			accountInfoCache.Add(account, info)
			return info, true
		}
	}
}

func CheckLogin(account, password string) (errs.AdminErrorCode, *models.Account) {
	dao := dbs.AccountDao{}
	defaultAccount, err := dao.FindByAccount("admin")
	if err != nil || defaultAccount == nil {
		//init default account
		err := dao.Create(dbs.AccountDao{
			Account:       "admin",
			Password:      tools.SHA1("123456"),
			CreatedTime:   time.Now(),
			UpdatedTime:   time.Now(),
			State:         0,
			ParentAccount: "",
		})
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}

	password = tools.SHA1(password)
	admin, err := dao.FindByAccountPassword(account, password)
	if err == nil && admin != nil {
		if admin.State != 0 {
			return errs.AdminErrorCode_AccountForbidden, nil
		}
		return errs.AdminErrorCode_Success, &models.Account{
			Account:       admin.Account,
			State:         admin.State,
			ParentAccount: admin.ParentAccount,
			// RoleId:        admin.RoleId,
			RoleType:    admin.RoleType,
			CreatedTime: admin.CreatedTime.UnixMilli(),
			UpdatedTime: admin.UpdatedTime.UnixMilli(),
		}
	} else if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_LoginFail, nil
}

func CheckAccountState(account string) errs.AdminErrorCode {
	dao := dbs.AccountDao{}
	admin, err := dao.FindByAccount(account)
	if err != nil || admin == nil {
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
		return errs.AdminErrorCode_AccountNotExist
	} else {
		if admin.State == 0 {
			return errs.AdminErrorCode_Success
		} else {
			return errs.AdminErrorCode_AccountForbidden
		}
	}
}

func UpdPassword(account, password, newPassword string) errs.AdminErrorCode {
	dao := dbs.AccountDao{}
	password = tools.SHA1(password)
	admin, err := dao.FindByAccountPassword(account, password)
	if err != nil || admin == nil {
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
		return errs.AdminErrorCode_UpdPwdFail
	}
	newPassword = tools.SHA1(newPassword)
	dao.UpdatePassword(account, newPassword)
	return errs.AdminErrorCode_Success
}

func AddAccount(ctx context.Context, account, password string, roleType int) errs.AdminErrorCode {
	parentAccount := ctxs.GetAccountFromCtx(ctx)
	curAccount, exist := GetAccountInfo(parentAccount)
	if !exist || curAccount == nil {
		return errs.AdminErrorCode_AccountNotExist
	}
	if curAccount.RoleType != RoleType_SuperAdmin {
		return errs.AdminErrorCode_NotPermission
	}
	dao := dbs.AccountDao{}
	password = tools.SHA1(password)
	err := dao.Create(dbs.AccountDao{
		Account:       account,
		Password:      password,
		ParentAccount: parentAccount,
		UpdatedTime:   time.Now(),
		CreatedTime:   time.Now(),
		RoleType:      roleType,
	})
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		return errs.AdminErrorCode_AccountExisted
	}
	return errs.AdminErrorCode_Success
}

func DisableAccounts(ctx context.Context, accounts []string, isDisable int) errs.AdminErrorCode {
	curAccount, exist := GetAccountInfo(ctxs.GetAccountFromCtx(ctx))
	if !exist || curAccount == nil {
		return errs.AdminErrorCode_AccountNotExist
	}
	if curAccount.RoleType != RoleType_SuperAdmin {
		return errs.AdminErrorCode_NotPermission
	}
	dao := dbs.AccountDao{}
	err := dao.UpdateState(accounts, isDisable)
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success
}

func DeleteAccounts(ctx context.Context, accounts []string) errs.AdminErrorCode {
	curAccount, exist := GetAccountInfo(ctxs.GetAccountFromCtx(ctx))
	if !exist || curAccount == nil {
		return errs.AdminErrorCode_AccountNotExist
	}
	if curAccount.RoleType != RoleType_SuperAdmin {
		return errs.AdminErrorCode_NotPermission
	}
	dao := dbs.AccountDao{}
	err := dao.DeleteAccounts(accounts)
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success
}

func BindApps(ctx context.Context, account string, appkeys []string) errs.AdminErrorCode {
	curAccount, exist := GetAccountInfo(ctxs.GetAccountFromCtx(ctx))
	if !exist || curAccount == nil {
		return errs.AdminErrorCode_AccountNotExist
	}
	if curAccount.RoleType != RoleType_SuperAdmin {
		return errs.AdminErrorCode_NotPermission
	}
	dao := dbs.AccountAppRelDao{}
	for _, appkey := range appkeys {
		dao.Create(dbs.AccountAppRelDao{
			AppKey:  appkey,
			Account: account,
		})
	}
	return errs.AdminErrorCode_Success
}

func UnBindApps(ctx context.Context, account string, appkeys []string) errs.AdminErrorCode {
	curAccount, exist := GetAccountInfo(ctxs.GetAccountFromCtx(ctx))
	if !exist || curAccount == nil {
		return errs.AdminErrorCode_AccountNotExist
	}
	if curAccount.RoleType != RoleType_SuperAdmin {
		return errs.AdminErrorCode_NotPermission
	}
	dao := dbs.AccountAppRelDao{}
	dao.BatchDelete(account, appkeys)
	return errs.AdminErrorCode_Success
}

func QryAccounts(ctx context.Context, limit int64, offset string) (errs.AdminErrorCode, *models.Accounts) {
	curAccount, exist := GetAccountInfo(ctxs.GetAccountFromCtx(ctx))
	if !exist || curAccount == nil {
		return errs.AdminErrorCode_AccountNotExist, nil
	}
	if curAccount.RoleType != RoleType_SuperAdmin {
		return errs.AdminErrorCode_NotPermission, nil
	}

	accounts := &models.Accounts{
		Items:   []*models.Account{},
		HasMore: false,
		Offset:  "",
	}
	dao := dbs.AccountDao{}
	offsetInt, err := tools.DecodeInt(offset)
	if err != nil {
		offsetInt = 0
	}
	dbAccounts, err := dao.QryAccounts(limit+1, offsetInt)
	if err == nil {
		if len(dbAccounts) > int(limit) {
			dbAccounts = dbAccounts[:len(dbAccounts)-1]
			accounts.HasMore = true
		}
		var id int64 = 0
		for _, dbAccount := range dbAccounts {
			accounts.Items = append(accounts.Items, &models.Account{
				Account:       dbAccount.Account,
				State:         dbAccount.State,
				CreatedTime:   dbAccount.CreatedTime.UnixMilli(),
				UpdatedTime:   dbAccount.UpdatedTime.UnixMilli(),
				ParentAccount: dbAccount.ParentAccount,
				// RoleId:        dbAccount.RoleId,
				RoleType: dbAccount.RoleType,
			})
			if dbAccount.ID > id {
				id = dbAccount.ID
			}
		}
		if id > 0 {
			offset, _ := tools.EncodeInt(id)
			accounts.Offset = offset
		}
	} else {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success, accounts
}
