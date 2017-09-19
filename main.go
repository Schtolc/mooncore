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
	// middleware
	e.Use(middleware.LoggerWithConfig(logger.Configure(conf.Logs.Access)))
	logfile := logger.OpenLogFile(conf.Logs.Main)

	// new output
	if err := syscall.Dup2(int(logfile.Fd()), int(os.Stderr.Fd())); err != nil {
		logrus.Fatal(err)
	}
	if err := syscall.Dup2(int(logfile.Fd()), int(os.Stdout.Fd())); err != nil {
		logrus.Fatal(err)
	}
	// init database
	db := database.Init(conf)
	defer db.Close()

	// new
	e.POST("/v1/register", handlers.Register(db))
	e.POST("/v1/login", handlers.Login(db))

	// groups
	AuthGroup := e.Group("/v1")
	jwtConfig := middleware.JWTConfig{
		SigningMethod: "HS256",
		Claims:     &handlers.JwtClaims{},
		SigningKey: []byte("secret"),
	}
	AuthGroup.Use(middleware.JWTWithConfig(jwtConfig))
	AuthGroup.Use(handlers.CheckJwtToken)
	AuthGroup.GET("/", handlers.Restricted)
	AuthGroup.GET("/ping", handlers.Ping)
	AuthGroup.GET("/ping_db", handlers.PingDb(db))


	logrus.Fatal(e.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}

