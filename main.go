package main

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/database"
	"github.com/Schtolc/mooncore/handlers"
	"github.com/Schtolc/mooncore/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	conf := config.Get()
	e.Use(middleware.LoggerWithConfig(logger.Configure(conf.Server.Logs.Access)))

	db := database.Init(conf)
	defer db.Close()

	e.GET("/ping", handlers.Ping)
	e.GET("/ping_db", handlers.PingDb(db))

	e.Logger.Fatal(e.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
