package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/dbcommons"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/routers"
)

func main() {
	//init configure
	if err := configures.InitConfigures(); err != nil {
		fmt.Println("Init Configures failed", err)
		return
	}
	//init log
	logs.InitLogs()
	//init mysql
	if err := dbcommons.InitMysql(); err != nil {
		fmt.Println("Init Mysql failed", err)
		return
	}
	//upgrade db
	dbcommons.Upgrade()

	httpServer := gin.Default()
	routers.Route(httpServer, "admingateway")
	routers.LoadJuggleChatAdminWeb(httpServer)
	go httpServer.Run(fmt.Sprintf(":%d", configures.Config.Port))

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		signal.Stop(sigChan)
		close(closeChan)
	}()

	<-closeChan
}
