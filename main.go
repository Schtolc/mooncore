package main

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/database"
	"github.com/Schtolc/mooncore/handlers"
	"github.com/Schtolc/mooncore/logger"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer logger.CatchError()
	e := echo.New()
	conf := config.Get()

	logger.Init(conf)
	e.Use(logger.Log(conf))

	db := database.Init(conf)
	defer db.Close()

	e.GET("/ping", handlers.Ping)
	e.GET("/ping_db", handlers.PingDb(db))

	log.Fatal(e.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
