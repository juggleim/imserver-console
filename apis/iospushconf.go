package apis

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/apnscredentials"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services"
	"github.com/juggleim/imserver-console/services/models"
)

// Validate middleware authenticates every route. Signature callers have no
// account; browser accounts use the existing role and app-binding checks.
func allowIosApp(ctx *gin.Context, appkey string) bool {
	account := ctx.GetString(string(ctxs.CtxKey_Account))
	if account == "" {
		return true
	}
	info, ok := services.GetAccountInfo(account)
	if ok && info.State == services.AccountState_Normal && (info.RoleType == services.RoleType_SuperAdmin || (dbs.AccountAppRelDao{}).CheckExist(appkey, account)) {
		return true
	}
	ctxs.FailHttpResp(ctx, errs.AdminErrorCode_NotPermission)
	return false
}

func iosPushItem(row *dbs.IosCertificateDao) *models.IosPushConfListItem {
	if row == nil {
		return nil
	}
	authType := row.AuthType
	if authType == "" {
		authType = "p12"
	}
	version := row.ConfigVersion
	if version < 1 {
		version = 1
	}
	return &models.IosPushConfListItem{
		AppKey: row.AppKey, Package: row.Package, IsProduct: row.IsProduct,
		CertPath: row.CertPath, VoipCertPath: row.VoipCertPath,
		AuthType: authType, P8KeyID: row.P8KeyID, P8TeamID: row.P8TeamID,
		P8KeyName: row.P8KeyName, HasP8Key: len(row.P8PrivateKey) > 0,
		ConfigVersion: version,
	}
}

func GetIosCer(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	if appkey == "" {
		failPushParam(ctx, "app_key is required")
		return
	}
	if !allowIosApp(ctx, appkey) {
		return
	}
	item, err := (dbs.IosCertificateDao{}).Find(appkey)
	if err != nil {
		failIosStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, iosPushItem(item))
}

func ListIosPushConfs(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	if appkey == "" {
		failPushParam(ctx, "app_key is required")
		return
	}
	if !allowIosApp(ctx, appkey) {
		return
	}
	rows, err := (dbs.IosCertificateDao{}).List(appkey)
	if err != nil {
		failIosStore(ctx, err)
		return
	}
	items := make([]*models.IosPushConfListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, iosPushItem(row))
	}
	ctxs.SuccessHttpResp(ctx, items)
}

type IosPushReq struct {
	AppKey          string `json:"app_key"`
	Package         string `json:"package"`
	OriginalPackage string `json:"original_package,omitempty"`
	IsProduct       int    `json:"is_product"`
	CertPath        string `json:"cert_path"`
	CertPwd         string `json:"cert_pwd"`
	VoipCertPath    string `json:"voip_cert_path"`
	VoipCertPwd     string `json:"voip_cert_pwd"`
	AuthType        string `json:"auth_type"`
	P8KeyID         string `json:"p8_key_id"`
	P8TeamID        string `json:"p8_team_id"`
	ConfigVersion   *int64 `json:"config_version"`
}

func SetIosPushConf(ctx *gin.Context) {
	var req IosPushReq
	if ctx.ShouldBindJSON(&req) != nil {
		failPushParam(ctx, "param illegal")
		return
	}
	req.AppKey, req.Package, req.OriginalPackage = strings.TrimSpace(req.AppKey), strings.TrimSpace(req.Package), strings.TrimSpace(req.OriginalPackage)
	if req.AppKey == "" || req.Package == "" || req.OriginalPackage == "" {
		failPushParam(ctx, "app_key, package and original_package are required")
		return
	}
	if !allowIosApp(ctx, req.AppKey) {
		return
	}
	item := dbs.IosCertificateDao{AppKey: req.AppKey, Package: req.Package, IsProduct: req.IsProduct,
		CertPwd: req.CertPwd, VoipCertPwd: req.VoipCertPwd, AuthType: req.AuthType,
		P8KeyID: req.P8KeyID, P8TeamID: req.P8TeamID, ExpectedConfigVersion: req.ConfigVersion}
	if err := (dbs.IosCertificateDao{}).Save(item, req.OriginalPackage); err != nil {
		failIosStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func UploadIosCer(ctx *gin.Context) {
	// Bound multipart parsing as well as individual reads; ordinary P12 files
	// retain ample room while P8 has the stricter shared 16 KiB bound.
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 8<<20)
	// Keep the bounded request in memory so multipart parsing never spills a
	// plaintext P8 key to a temporary file after larger P12 parts.
	if err := ctx.Request.ParseMultipartForm(8 << 20); err != nil {
		failPushParam(ctx, "invalid or oversized upload")
		return
	}
	defer ctx.Request.MultipartForm.RemoveAll()
	appkey, packageName := strings.TrimSpace(ctx.PostForm("app_key")), strings.TrimSpace(ctx.PostForm("package"))
	if appkey == "" || packageName == "" {
		failPushParam(ctx, "app_key and package are required")
		return
	}
	if !allowIosApp(ctx, appkey) {
		return
	}
	isProduct := 0
	if value := ctx.PostForm("is_product"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || (parsed != 0 && parsed != 1) {
			failPushParam(ctx, "invalid certificate environment")
			return
		}
		isProduct = parsed
	}
	item := dbs.IosCertificateDao{AppKey: appkey, Package: packageName, IsProduct: isProduct,
		CertPwd: ctx.PostForm("cert_pwd"), VoipCertPwd: ctx.PostForm("voip_cert_pwd"),
		AuthType: ctx.PostForm("auth_type"), P8KeyID: ctx.PostForm("p8_key_id"), P8TeamID: ctx.PostForm("p8_team_id")}
	var privateKey []byte
	if values, present := ctx.Request.MultipartForm.Value["config_version"]; present {
		value := values[0]
		version, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			failPushParam(ctx, "invalid config_version")
			return
		}
		item.ExpectedConfigVersion = &version
	}
	defer func() { clear(privateKey) }()
	for _, upload := range []struct {
		field string
		data  *[]byte
		name  *string
		limit int64
	}{
		{"ioscer", &item.Certificate, &item.CertPath, 8 << 20},
		{"voip_ioscer", &item.VoipCert, &item.VoipCertPath, 8 << 20},
		{"p8_file", &privateKey, &item.P8KeyName, apnscredentials.MaxPrivateKeySize},
	} {
		header, err := ctx.FormFile(upload.field)
		if errors.Is(err, http.ErrMissingFile) {
			continue
		}
		if err != nil || header.Size <= 0 || header.Size > upload.limit || len(header.Filename) > 255 {
			failPushParam(ctx, "invalid credential file or size")
			return
		}
		file, err := header.Open()
		if err != nil {
			failPushParam(ctx, "unable to open credential file")
			return
		}
		data, readErr := io.ReadAll(io.LimitReader(file, upload.limit+1))
		file.Close()
		if readErr != nil || len(data) == 0 || int64(len(data)) > upload.limit {
			clear(data)
			failPushParam(ctx, "invalid credential file or size")
			return
		}
		*upload.data, *upload.name = data, header.Filename
	}
	if err := (dbs.IosCertificateDao{}).Save(item, ctx.PostForm("original_package"), privateKey); err != nil {
		failIosStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func failIosStore(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, dbs.ErrIosVersionConflict):
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "iOS configuration changed; reload before saving", "data": nil})
	case errors.Is(err, dbs.ErrIosCredentials):
		failPushParam(ctx, "invalid iOS credentials or metadata; for P8 verify the key and IDs")
	case errors.Is(err, dbs.ErrPushConfConflict), errors.Is(err, dbs.ErrPushConfNotFound):
		failPushStore(ctx, err)
	default:
		// Database errors can contain SQL parameters. Never forward or log them.
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ServerErr, "iOS push configuration operation failed")
	}
}
