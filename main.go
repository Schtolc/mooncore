package main

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/database"
	"github.com/Schtolc/mooncore/logger"
	"github.com/Schtolc/mooncore/server"
	"github.com/sirupsen/logrus"
)

func main() {
	conf := config.Get()

	logger.Init(conf)

	db := database.Init(conf)
	defer db.Close()

	server := server.Init(conf, db)
	logrus.Fatal(server.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
