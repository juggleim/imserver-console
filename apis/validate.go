package apis

import (
	"fmt"

	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"
)

const (
	Header_RequestId     string = "request-id"
	Header_Authorization string = "Authorization"
	Header_Timestamp     string = "timestamp"
	Header_Nonce         string = "timestamp"
	Header_Signature     string = "signature"
)

func Validate(ctx *gin.Context) {
	session := fmt.Sprintf("admin_%s", tools.GenerateUUIDShort11())
	ctx.Header(Header_RequestId, session)
	ctx.Set(string(ctxs.CtxKey_Session), session)

	urlPath := ctx.Request.URL.Path
	if strings.HasSuffix(urlPath, "/login") {
		return
	}

	signature := ctx.Request.Header.Get(Header_Signature)
	if signature != "" && configures.Config.AdminSecret != "" {
		secret := configures.Config.AdminSecret
		nonce := ctx.Request.Header.Get(Header_Nonce)
		tsStr := ctx.Request.Header.Get(Header_Timestamp)
		if nonce == "" || tsStr == "" {
			ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
				Code: errs.AdminErrorCode_ParamError,
			})
			ctx.Abort()
			return
		}
		str := fmt.Sprintf("%s%s%s", secret, nonce, tsStr)
		sig := tools.SHA1(str)
		if sig != signature {
			ctx.JSON(http.StatusUnauthorized, &ctxs.ApiErrorMsg{
				Code: errs.AdminErrorCode_AuthFail,
				Msg:  "auth failed",
			})
			ctx.Abort()
			return
		}
	} else {
		authStr := ctx.Request.Header.Get(Header_Authorization)
		account, err := validateAuthorization(authStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, &ctxs.ApiErrorMsg{
				Code: errs.AdminErrorCode_AuthFail,
				Msg:  "auth failed",
			})
			ctx.Abort()
			return
		}
		//check account
		code := services.CheckAccountState(account)
		if code != errs.AdminErrorCode_Success {
			ctx.JSON(http.StatusUnauthorized, &ctxs.ApiErrorMsg{
				Code: code,
				Msg:  "auth failed",
			})
			ctx.Abort()
			return
		}
		ctx.Set(string(ctxs.CtxKey_Account), account)
	}
}

func GetLoginedAccount(ctx *gin.Context) string {
	if account, ok := ctx.Value(ctxs.CtxKey_Account).(string); ok {
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
