package main

import (
	"github.com/Schtolc/mooncore/dependencies"
	//"github.com/pkg/profile"
	"github.com/sirupsen/logrus"
)

func main() {
	//defer profile.Start().Stop()
	conf := dependencies.ConfigInstance()
	InitLogs(conf)

	db := dependencies.DBInstance()
	defer db.Close()

	server := InitServer(conf)
	logrus.Fatal(server.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
