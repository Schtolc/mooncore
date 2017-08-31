package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"mooncore/cfg"
	"mooncore/handlers"
	"mooncore/logger"
)

func main() {
	e := echo.New()

	conf := cfg.GetAppConfig("app.conf")
	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(logger.Configure(conf.Log.Access)))

	db := database.InitDB()

	defer db.Close()

	e.GET("/ping", handlers.Ping)
	e.GET("/ping_db", handlers.PingDb(db))
	e.POST("/create_user", handlers.SaveUser(db))

	e.Logger.Fatal(e.Start(conf.Hostbase.Host + ":" + conf.Hostbase.Port))
}
