package apis

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

func GetIosCer(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	if appkey == "" {
		failPushParam(ctx, "app_key is required")
		return
	}
	item, err := (dbs.IosCertificateDao{}).Find(appkey)
	if err != nil {
		failPushStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, item)
}

func ListIosPushConfs(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	if appkey == "" {
		failPushParam(ctx, "app_key is required")
		return
	}
	rows, err := (dbs.IosCertificateDao{}).List(appkey)
	if err != nil {
		failPushStore(ctx, err)
		return
	}
	items := make([]models.IosPushConfListItem, 0, len(rows))
	for _, row := range rows {
		item := models.IosPushConfListItem{
			AppKey:       row.AppKey,
			Package:      row.Package,
			IsProduct:    row.IsProduct,
			CertPath:     row.CertPath,
			CertPwd:      row.CertPwd,
			VoipCertPwd:  row.VoipCertPwd,
			VoipCertPath: row.VoipCertPath,
		}
		items = append(items, item)
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
}

func SetIosPushConf(ctx *gin.Context) {
	var req IosPushReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		failPushParam(ctx, "param illegal")
		return
	}
	req.AppKey = strings.TrimSpace(req.AppKey)
	req.Package = strings.TrimSpace(req.Package)
	req.OriginalPackage = strings.TrimSpace(req.OriginalPackage)
	if req.AppKey == "" || req.Package == "" || (req.IsProduct != 0 && req.IsProduct != 1) {
		failPushParam(ctx, "app_key, package and valid certificate environment are required")
		return
	}
	if req.OriginalPackage == "" {
		failPushParam(ctx, "a certificate file is required for a new iOS configuration")
		return
	}

	dao := dbs.IosCertificateDao{}
	existing, err := dao.FindByPackage(req.AppKey, req.OriginalPackage)
	if err != nil {
		failPushStore(ctx, err)
		return
	}
	item := *existing
	item.Package = req.Package
	item.IsProduct = req.IsProduct
	if req.CertPwd != "" {
		item.CertPwd = req.CertPwd
	}
	if req.VoipCertPwd != "" {
		item.VoipCertPwd = req.VoipCertPwd
	}
	if err := dao.Save(item, req.OriginalPackage); err != nil {
		failPushStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func UploadIosCer(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.PostForm("app_key"))
	packageName := strings.TrimSpace(ctx.PostForm("package"))
	originalPackage := strings.TrimSpace(ctx.PostForm("original_package"))
	if appkey == "" || packageName == "" {
		failPushParam(ctx, "app_key and package are required")
		return
	}

	isProduct := 0
	if value := strings.TrimSpace(ctx.PostForm("is_product")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || (parsed != 0 && parsed != 1) {
			failPushParam(ctx, "invalid certificate environment")
			return
		}
		isProduct = parsed
	}

	dao := dbs.IosCertificateDao{}
	item := dbs.IosCertificateDao{AppKey: appkey, Package: packageName, IsProduct: isProduct}
	if originalPackage != "" {
		existing, err := dao.FindByPackage(appkey, originalPackage)
		if err != nil {
			failPushStore(ctx, err)
			return
		}
		item = *existing
		item.Package = packageName
		item.IsProduct = isProduct
	}
	if value := ctx.PostForm("cert_pwd"); value != "" {
		item.CertPwd = value
	}
	if value := ctx.PostForm("voip_cert_pwd"); value != "" {
		item.VoipCertPwd = value
	}

	fileHeader, fileErr := ctx.FormFile("ioscer")
	if fileErr == nil {
		file, err := fileHeader.Open()
		if err != nil {
			failPushParam(ctx, "unable to open iOS certificate")
			return
		}
		defer file.Close()
		item.Certificate, err = io.ReadAll(file)
		if err != nil {
			failPushParam(ctx, "unable to read iOS certificate")
			return
		}
		item.CertPath = fileHeader.Filename
	} else if !errors.Is(fileErr, http.ErrMissingFile) {
		failPushParam(ctx, "invalid iOS certificate")
		return
	}

	voipHeader, voipErr := ctx.FormFile("voip_ioscer")
	if voipErr == nil {
		file, err := voipHeader.Open()
		if err != nil {
			failPushParam(ctx, "unable to open VoIP certificate")
			return
		}
		defer file.Close()
		item.VoipCert, err = io.ReadAll(file)
		if err != nil {
			failPushParam(ctx, "unable to read VoIP certificate")
			return
		}
		item.VoipCertPath = voipHeader.Filename
	} else if !errors.Is(voipErr, http.ErrMissingFile) {
		failPushParam(ctx, "invalid VoIP certificate")
		return
	}

	if len(item.Certificate) == 0 || item.CertPath == "" || item.CertPwd == "" {
		failPushParam(ctx, "certificate file and password are required")
		return
	}
	if len(item.VoipCert) > 0 && item.VoipCertPwd == "" {
		failPushParam(ctx, "VoIP certificate password is required")
		return
	}
	if err := dao.Save(item, originalPackage); err != nil {
		failPushStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}
