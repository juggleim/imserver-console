package routers

import (
	"embed"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Vite 会产出以 "_" 开头的 chunk（如 _plugin-vue_export-helper-*.js），
// 需使用 all: 前缀，否则 go:embed 默认会排除这类文件。
//
//go:embed all:admin
var adminFiles embed.FS

func LoadJuggleChatAdminWeb(httpServer *gin.Engine) {
	httpServer.GET("/assets/*filepath", serveAsset)
	httpServer.GET("/", dashboardPage)
	httpServer.GET("/login", dashboardPage)
	httpServer.GET("/dashboard", dashboardPage)
	httpServer.NoRoute(spaFallback)
}

func dashboardPage(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK, readTextFile("admin/index.html"))
}

var htmlCache sync.Map

var assetCache sync.Map

type assetCacheItem struct {
	contentType string
	body        []byte
}

func serveAsset(ctx *gin.Context) {
	filePath := strings.TrimPrefix(ctx.Param("filepath"), "/")
	if filePath == "" || strings.Contains(filePath, "..") {
		ctx.Status(http.StatusNotFound)
		return
	}

	cacheKey := "/assets/" + filePath
	if cached, ok := assetCache.Load(cacheKey); ok {
		item := cached.(assetCacheItem)
		ctx.Data(http.StatusOK, item.contentType, item.body)
		return
	}

	embedPath := "admin/assets/" + filePath
	body, err := adminFiles.ReadFile(embedPath)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	contentType := mime.TypeByExtension(path.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	item := assetCacheItem{contentType: contentType, body: body}
	assetCache.Store(cacheKey, item)
	ctx.Data(http.StatusOK, contentType, body)
}

func spaFallback(ctx *gin.Context) {
	requestPath := ctx.Request.URL.Path
	if strings.HasPrefix(requestPath, "/admingateway") || strings.HasPrefix(requestPath, "/assets/") || path.Ext(requestPath) != "" {
		ctx.Status(http.StatusNotFound)
		return
	}
	dashboardPage(ctx)
}

func readTextFile(embedPath string) string {
	if cached, ok := htmlCache.Load(embedPath); ok {
		return cached.(string)
	}
	bs, err := adminFiles.ReadFile(embedPath)
	if err != nil {
		fmt.Println("read file failed:", embedPath, err)
		return ""
	}
	body := string(bs)
	htmlCache.Store(embedPath, body)
	return body
}
