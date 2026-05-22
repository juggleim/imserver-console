package apis

import (
	"bufio"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/imsdk"
	"github.com/juggleim/imserver-console/commons/tools"
	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

type QrySensitiveWordsResp struct {
	Items      []*SensitiveWord `json:"items"`
	IsFinished bool             `json:"is_finished"`
	Total      int32            `json:"total"`
}

type SensitiveWord struct {
	AppKey   string `json:"app_key,omitempty"`
	Id       string `json:"id,omitempty"`
	Word     string `json:"word"`
	WordType int    `json:"word_type"`
}

func getSensitiveSdk(ctx *gin.Context, appKey string) (*juggleimsdk.JuggleIMSdk, bool) {
	if appKey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "app_key required")
		return nil, false
	}
	sdk := imsdk.GetImSdk(appKey)
	if sdk == nil {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "app not found")
		return nil, false
	}
	return sdk, true
}

func handleSdkErr(ctx *gin.Context, code juggleimsdk.ApiCode, err error) bool {
	if err != nil {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ServerErr, err.Error())
		return true
	}
	if code != juggleimsdk.ApiCode_Success {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode(code), "")
		return true
	}
	return false
}

func SensitiveWords(ctx *gin.Context) {
	sizeStr := ctx.Query("size")
	var size int64 = 50
	if sizeStr != "" {
		intVal, err := tools.String2Int64(sizeStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			size = intVal
		}
	}
	pageStr := ctx.Query("page")
	var page int64 = 1
	if pageStr != "" {
		intVal, err := tools.String2Int64(pageStr)
		if err == nil && intVal > 0 {
			page = intVal
		}
	}
	word := ctx.Query("word")
	wordTypeStr := ctx.Query("word_type")
	var wordType int
	if wordTypeStr != "" {
		intVal, err := tools.String2Int64(wordTypeStr)
		if err == nil {
			wordType = int(intVal)
		}
	}

	appkey := ctx.Query("app_key")
	sdk, ok := getSensitiveSdk(ctx, appkey)
	if !ok {
		return
	}

	resp, code, _, err := sdk.QrySensitiveWords(int(size), int(page), word, wordType)
	if handleSdkErr(ctx, code, err) {
		return
	}

	res := &QrySensitiveWordsResp{
		Items:      []*SensitiveWord{},
		IsFinished: resp.IsFinished,
		Total:      resp.Total,
	}
	for _, senWord := range resp.Items {
		res.Items = append(res.Items, &SensitiveWord{
			Id:       senWord.Id,
			Word:     senWord.Word,
			WordType: senWord.WordType,
		})
	}
	ctxs.SuccessHttpResp(ctx, res)
}

func ImportSensitiveWords(ctx *gin.Context) {
	appKey := ctxs.GetAppKeyFromCtx(ctxs.ToCtx(ctx))
	if appKey == "" {
		appKey = ctx.Query("app_key")
	}
	if appKey == "" {
		appKey = ctx.PostForm("app_key")
	}
	sdk, ok := getSensitiveSdk(ctx, appKey)
	if !ok {
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, err.Error())
		return
	}
	f, err := file.Open()
	if err != nil {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ServerErr, err.Error())
		return
	}
	defer f.Close()

	var items []*juggleimsdk.SensitiveWord
	reader := bufio.NewReader(f)
	for {
		bs, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ServerErr, err.Error())
			return
		}
		word := strings.TrimSpace(string(bs))
		if word == "" {
			continue
		}
		items = append(items, &juggleimsdk.SensitiveWord{
			Word:     word,
			WordType: 2,
		})
	}

	code, _, err := sdk.AddSensitiveWords(juggleimsdk.SensitiveWords{Items: items})
	if handleSdkErr(ctx, code, err) {
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func AddSensitiveWord(ctx *gin.Context) {
	var req SensitiveWord
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Word == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "param illegal")
		return
	}

	sdk, ok := getSensitiveSdk(ctx, req.AppKey)
	if !ok {
		return
	}

	code, _, err := sdk.AddSensitiveWords(juggleimsdk.SensitiveWords{
		Items: []*juggleimsdk.SensitiveWord{
			{
				Word:     req.Word,
				WordType: req.WordType,
			},
		},
	})
	if handleSdkErr(ctx, code, err) {
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func DeleteSensitiveWord(ctx *gin.Context) {
	var req SensitiveWord
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Word == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "param illegal")
		return
	}

	sdk, ok := getSensitiveSdk(ctx, req.AppKey)
	if !ok {
		return
	}

	code, _, err := sdk.DeleteSensitiveWords(juggleimsdk.DelSensitiveWordsReq{
		Words: []string{req.Word},
	})
	if handleSdkErr(ctx, code, err) {
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}
