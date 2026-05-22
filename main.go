package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/dbcommons"
	"github.com/juggleim/imserver-console/commons/logs"
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
	httpServer.Run(fmt.Sprintf(":%d", configures.Config.Port))
}
