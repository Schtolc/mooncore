package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"mooncore/config"
	"mooncore/database"
	"mooncore/handlers"
	"mooncore/logger"
)

func main() {
	e := echo.New()

	conf := config.Get()
	e.Use(middleware.LoggerWithConfig(logger.Configure(conf.Server.Logs.Access)))
	e.Use(middleware.RequestID())

	db := database.InitDB(conf)
	defer db.Close()

	e.POST("/ping", handlers.EchoPing)
	e.GET("/ping_db", handlers.PingDb(db))

	e.Logger.Fatal(e.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
