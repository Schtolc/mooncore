package main

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/database"
	"github.com/sirupsen/logrus"
)

func main() {
	conf := config.Instance()
	InitLogs(conf)

	db := database.Instance()
	defer db.Close()

	server := InitServer(conf, db)
	logrus.Fatal(server.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
