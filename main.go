package main

import (
	"github.com/Schtolc/mooncore/database"
	"github.com/Schtolc/mooncore/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	conf := utils.GetConfig()
	InitLogs(conf)

	db := database.InitDatabase(conf)
	defer db.Close()

	server := InitServer(conf, db)
	logrus.Fatal(server.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
