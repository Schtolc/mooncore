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

	mainLog := logger.OpenLogFile(conf.Logs.Main)
	if err := syscall.Dup2(int(mainLog.Fd()), int(os.Stderr.Fd())); err != nil {
		logrus.Fatal(err)
	}
	if err := syscall.Dup2(int(mainLog.Fd()), int(os.Stdout.Fd())); err != nil {
		logrus.Fatal(err)
	}

	db := database.Init(conf)
	defer db.Close()

	h := handlers.Init(db)

	e.POST("/v1/sign_up", h.SignUp)
	e.POST("/v1/sign_in", h.SignIn)

	AuthGroup := e.Group("/v1")
	jwtConfig := middleware.JWTConfig{
		SigningMethod: "HS256",
		Claims:        &handlers.JwtClaims{},
		SigningKey:    handlers.SigningKey,
	}
	AuthGroup.Use(middleware.JWTWithConfig(jwtConfig))
	AuthGroup.Use(h.CheckJwtToken)
	AuthGroup.GET("/ping", h.Ping)
	AuthGroup.GET("/ping_db", h.PingDb)

	logrus.Fatal(e.Start(conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port))
}
