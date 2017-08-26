package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"mooncore/handlers"
	"mooncore/logger"
)


func main() {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(logger.Configure("logfile")))

	db := ConnectDB()
	defer db.Close()

	e.GET("/ping", handlers.Ping)
	e.GET("/ping_db", handlers.PingDb(db))

	e.Logger.Fatal(e.Start(":1323"))
}



