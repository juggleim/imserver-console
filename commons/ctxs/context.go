package ctxs

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/errs"
)

type CtxKey string

const (
	CtxKey_AppKey      CtxKey = "CtxKey_AppKey"
	CtxKey_RequesterId CtxKey = "CtxKey_RequesterId"
	CtxKey_Session     CtxKey = "CtxKey_Session"
	CtxKey_Account     CtxKey = "CtxKey_Account"
	CtxKey_RoleType    CtxKey = "CtxKey_RoleType"
)

const (
	Header_RequestId     string = "request-id"
	Header_Authorization string = "Authorization"
	Header_Timestamp     string = "timestamp"
	Header_Nonce         string = "timestamp"
	Header_Signature     string = "signature"
)

func GetLoginedAccount(ctx *gin.Context) string {
	if account, ok := ctx.Value(CtxKey_Account).(string); ok {
		return account
	}
	return ""
}

var defaultJwtkey = []byte("jug9le1m")

func getJwtKey() []byte {
	adminSecret := configures.Config.AdminSecret
	if adminSecret != "" {
		return []byte(adminSecret)
	} else {
		return defaultJwtkey
	}
}

type Claims struct {
	Account string
	jwt.RegisteredClaims
}

func generateAuthorization(account string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Account: account,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expireTime,
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			Issuer:  "aabbcc",
			Subject: "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJwtKey())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthTest(auth string) (string, error) {
	return validateAuthorization(auth)
}

func validateAuthorization(authorization string) (string, error) {
	token, claims, err := parseToken(authorization)
	if err != nil || !token.Valid {
		return "", fmt.Errorf("auth fail")
	}
	return claims.Account, nil
}

func parseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return getJwtKey(), nil
	})
	return token, Claims, err
}

func ToCtx(ginCtx *gin.Context) context.Context {
	rpcCtx := context.Background()
	appkey := ginCtx.GetString(string(CtxKey_AppKey))
	if appkey != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_AppKey, appkey)
	}
	session := ginCtx.GetString(string(CtxKey_Session))
	if session != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_Session, session)
	}
	currentUserId := ginCtx.GetString(string(CtxKey_RequesterId))
	if currentUserId != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_RequesterId, currentUserId)
	}
	account := ginCtx.GetString(string(CtxKey_Account))
	if account != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_Account, account)
	}
	return rpcCtx
}

func ToCtxWithRequester(ginCtx *gin.Context, requestId string) context.Context {
	rpcCtx := ToCtx(ginCtx)
	if requestId != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_RequesterId, requestId)
	}
	return rpcCtx
}

func GetAppKeyFromCtx(ctx context.Context) string {
	if appKey, ok := ctx.Value(CtxKey_AppKey).(string); ok {
		return appKey
	}
	return ""
}

func SetAppKeyToCtx(ctx context.Context, appkey string) context.Context {
	if appkey != "" {
		newCtx := context.WithValue(ctx, CtxKey_AppKey, appkey)
		return newCtx
	}
	return ctx
}

func GetRequesterIdFromCtx(ctx context.Context) string {
	if requesterId, ok := ctx.Value(CtxKey_RequesterId).(string); ok {
		return requesterId
	}
	return ""
}

func GetSessionFromCtx(ctx context.Context) string {
	if session, ok := ctx.Value(CtxKey_Session).(string); ok {
		return session
	}
	return ""
}

func GetAccountFromCtx(ctx context.Context) string {
	if requesterId, ok := ctx.Value(CtxKey_Account).(string); ok {
		return requesterId
	}
	return ""
}

type ApiErrorMsg struct {
	HttpCode int                 `json:"-"`
	Code     errs.AdminErrorCode `json:"code"`
	Msg      string              `json:"msg"`
}

type SuccHttpResp struct {
	ApiErrorMsg
	Data interface{} `json:"data"`
}

func SuccessHttpResp(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, SuccHttpResp{
		ApiErrorMsg: ApiErrorMsg{
			Code: 0,
			Msg:  "success",
		},
		Data: data,
	})
}

func FailHttpResp(ctx *gin.Context, code errs.AdminErrorCode, msgs ...string) {
	var msg = "fail"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	ctx.JSON(http.StatusOK, SuccHttpResp{
		ApiErrorMsg: ApiErrorMsg{
			Code: code,
			Msg:  msg,
		},
	})
}
