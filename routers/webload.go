package routers

import (
	"embed"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

//go:embed admin/*
var adminFiles embed.FS

func LoadJuggleChatAdminWeb(httpServer *gin.Engine) {
	files, err := adminFiles.ReadDir("admin/assets")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		if !f.IsDir() {
			httpServer.GET("/assets/"+f.Name(), assetsFile)
		}
	}
	httpServer.GET("/", dashboardPage)
	httpServer.GET("login", dashboardPage)
	httpServer.GET("/dashboard", dashboardPage)
}

func dashboardPage(ctx *gin.Context) {
	ctx.Writer.Header().Add("Content-Type", "text/html; charset=utf-8")

	var body string
	cacheBody, ok := htmlCache.Load("index.html")
	if ok {
		body = cacheBody.(string)
	} else {
		body = ReadFromFile("admin/index.html")
		htmlCache.Store("index.html", body)
	}
	ctx.String(200, body)
}

var htmlCache sync.Map

func assetsFile(ctx *gin.Context) {
	filePath := ctx.Request.URL.Path
	if strings.HasSuffix(filePath, ".js") {
		ctx.Writer.Header().Add("Content-Type", "application/javascript")
	} else if strings.HasSuffix(filePath, ".css") {
		ctx.Writer.Header().Add("Content-Type", "text/css")
	} else if strings.HasSuffix(filePath, ".png") {
		ctx.Writer.Header().Add("Content-Type", "image/png")
	} else if strings.HasSuffix(filePath, ".ico") {
		ctx.Writer.Header().Add("Content-Type", "image/x-icon")
	}
	var body string
	if cacheBody, ok := htmlCache.Load(filePath); ok {
		body = cacheBody.(string)
	} else {
		body = ReadFromFile("admin" + filePath)
		htmlCache.Store(filePath, body)
	}
	ctx.String(200, body)
}

func ReadFromFile(path string) string {
	// bs, err := os.ReadFile(path)
	bs, err := adminFiles.ReadFile(path)
	if err != nil {
		fmt.Println("read file failed:", err)
		return ""
	}
	return string(bs)
}
