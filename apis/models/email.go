package models

import servicemodels "github.com/juggleim/imserver-console/services/models"

type EmailConf struct {
	AppKey string                        `json:"app_key"`
	Conf   *servicemodels.MailEngineConf `json:"conf"`
}
