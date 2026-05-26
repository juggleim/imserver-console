package services

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	apimodels "github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/caches"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/commons/fileengines"
)

type FileCredErr int

const (
	FileCredSuccess FileCredErr = 0
	FileCredNoOSS   FileCredErr = 1
	FileCredSignErr FileCredErr = 2
)

func GetFileCred(ctx context.Context, req *apimodels.QryFileCredReq) (FileCredErr, *apimodels.QryFileCredResp) {
	appkey := req.AppKey
	fileConf := getFileConf(ctx, appkey)
	if fileConf == nil || fileConf == notExistFileConf {
		return FileCredNoOSS, nil
	}
	dir := fileTypeToDir(req.FileType)
	switch fileConf.FileEngine {
	case fileengines.ChannelQiNiu:
		if fileConf.QiNiu == nil {
			return FileCredNoOSS, nil
		}
		uploadToken, domain := fileConf.QiNiu.UploadToken(req.Ext)
		return FileCredSuccess, &apimodels.QryFileCredResp{
			OssType: apimodels.OssType_QiNiu,
			QiNiuCredResp: &apimodels.QiNiuCredResp{
				Domain: domain,
				Token:  uploadToken,
			},
		}
	case fileengines.ChannelMinio:
		if fileConf.Minio == nil {
			return FileCredNoOSS, nil
		}
		signedURL, err := fileConf.Minio.PreSignedURL(req.Ext, dir)
		if err != nil {
			return FileCredSignErr, nil
		}
		return FileCredSuccess, &apimodels.QryFileCredResp{
			OssType: apimodels.OssType_Minio,
			PreSignResp: &apimodels.PreSignResp{
				Url: signedURL,
			},
		}
	case fileengines.ChannelAws:
		if fileConf.S3 == nil {
			return FileCredNoOSS, nil
		}
		signedURL, err := fileConf.S3.PreSignedURL(req.Ext, dir)
		if err != nil {
			return FileCredSignErr, nil
		}
		return FileCredSuccess, &apimodels.QryFileCredResp{
			OssType: apimodels.OssType_S3,
			PreSignResp: &apimodels.PreSignResp{
				Url: signedURL,
			},
		}
	case fileengines.ChannelOss:
		if fileConf.Oss == nil {
			return FileCredNoOSS, nil
		}
		signedURL, err := fileConf.Oss.PreSignedURL(req.Ext, dir)
		if err != nil {
			return FileCredSignErr, nil
		}
		resp := fileConf.Oss.PostSign(req.Ext, dir)
		return FileCredSuccess, &apimodels.QryFileCredResp{
			OssType: apimodels.OssType_Oss,
			PreSignResp: &apimodels.PreSignResp{
				Url:         signedURL,
				ObjKey:      resp.ObjKey,
				Policy:      resp.Policy,
				SignVersion: resp.SignVersion,
				Credential:  resp.Credential,
				Date:        resp.Date,
				Signature:   resp.Signature,
			},
		}
	default:
		return FileCredNoOSS, nil
	}
}

type FileConfItem struct {
	AppKey     string
	FileEngine string
	QiNiu      *fileengines.QiNiuStorage
	Oss        *fileengines.OssStorage
	Minio      *fileengines.MinioStorage
	S3         *fileengines.S3Storage
}

var fileConfCache *caches.LruCache
var fileLock *sync.RWMutex
var notExistFileConf *FileConfItem

func init() {
	fileConfCache = caches.NewLruCacheWithAddReadTimeout("fileconf_caches", 1000, nil, 5*time.Minute, 10*time.Minute)
	fileLock = &sync.RWMutex{}
	notExistFileConf = &FileConfItem{}
}

func getFileConf(ctx context.Context, appkey string) *FileConfItem {
	if obj, exist := fileConfCache.Get(appkey); exist {
		return obj.(*FileConfItem)
	}
	fileLock.Lock()
	defer fileLock.Unlock()
	if obj, exist := fileConfCache.Get(appkey); exist {
		return obj.(*FileConfItem)
	}
	fileConf, err := loadFileConfFromDb(appkey)
	if err != nil {
		fileConf = notExistFileConf
	}
	fileConfCache.Add(appkey, fileConf)
	return fileConf
}

func loadFileConfFromDb(appkey string) (*FileConfItem, error) {
	dao := dbs.FileConfDao{}
	conf, err := dao.FindEnableFileConf(appkey)
	if err != nil {
		return nil, err
	}
	fileConf := &FileConfItem{
		AppKey:     appkey,
		FileEngine: conf.Channel,
	}
	var confData = make(map[string]interface{})
	_ = json.Unmarshal([]byte(conf.Conf), &confData)

	switch conf.Channel {
	case fileengines.ChannelQiNiu:
		c := tools.MapToStruct[fileengines.QiNiuConfig](confData)
		fileConf.QiNiu = fileengines.NewQiNiu(c)
	case fileengines.ChannelMinio:
		c := tools.MapToStruct[fileengines.MinioConfig](confData)
		fileConf.Minio = fileengines.NewMinio(c)
	case fileengines.ChannelOss:
		c := tools.MapToStruct[fileengines.OssConfig](confData)
		fileConf.Oss = fileengines.NewOss(c)
	case fileengines.ChannelAws:
		c := tools.MapToStruct[fileengines.S3Config](confData)
		fileConf.S3 = fileengines.NewS3Storage(fileengines.WithConf(c))
	}
	return fileConf, nil
}

func fileTypeToDir(fileType apimodels.FileType) string {
	switch fileType {
	case apimodels.FileType_Image:
		return "images"
	case apimodels.FileType_Video:
		return "videos"
	case apimodels.FileType_Audio:
		return "audios"
	case apimodels.FileType_File:
		return "files"
	case apimodels.FileType_Log:
		return "logs"
	default:
		return "files"
	}
}
