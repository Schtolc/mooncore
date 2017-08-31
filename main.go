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

	conf := config.GetAppConfig("app.conf")
	e.Use(middleware.LoggerWithConfig(logger.Configure(conf.Log.Access)))
	e.Use(middleware.RequestID())

	db := database.InitDB(conf)
	defer db.Close()

	e.GET("/ping", handlers.Ping)
	e.GET("/ping_db", handlers.PingDb(db))

	e.Logger.Fatal(e.Start(conf.Hostbase.Host + ":" + conf.Hostbase.Port))
}
