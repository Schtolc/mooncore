package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	conf := GetConfig()
	InitLogs(conf)

	db := InitDatabase(conf)
	defer db.Close()

	server := InitServer(conf, db)
	logrus.Fatal(server.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
