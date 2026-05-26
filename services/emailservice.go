package services

import (
	"context"

	apimodels "github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	servicemodels "github.com/juggleim/imserver-console/services/models"
)

func SetEmailConf(ctx context.Context, req *apimodels.EmailConf) errs.AdminErrorCode {
	dao := dbs.AppExtDao{}
	_ = dao.Upsert(req.AppKey, "mail_engine_conf", tools.ToJson(req.Conf))
	return errs.AdminErrorCode_Success
}

func GetEmailConf(ctx context.Context, appkey string) (errs.AdminErrorCode, *servicemodels.MailEngineConf) {
	emailConf := &servicemodels.MailEngineConf{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "mail_engine_conf")
	if err == nil && conf != nil {
		_ = tools.JsonUnMarshal([]byte(conf.AppItemValue), emailConf)
	}
	return errs.AdminErrorCode_Success, emailConf
}
