package main

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/database"
	"github.com/Schtolc/mooncore/handlers"
	"github.com/Schtolc/mooncore/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func main() {
	e := echo.New()
	conf := config.Get()

	e.Use(middleware.LoggerWithConfig(logger.Configure(conf.Logs.Access)))
	logfile := logger.OpenLogFile(conf.Logs.Main)

	if err := syscall.Dup2(int(logfile.Fd()), int(os.Stderr.Fd())); err != nil {
		logrus.Fatal(err)
	}
	if err := syscall.Dup2(int(logfile.Fd()), int(os.Stdout.Fd())); err != nil {
		logrus.Fatal(err)
	}

	db := database.Init(conf)
	defer db.Close()

	e.GET("/ping", handlers.Ping)
	e.GET("/ping_db", handlers.PingDb(db))

	logrus.Fatal(e.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
