package main

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"os"
)

// InitServer inits echo server: sets access logs and handlers
func InitServer(config *dependencies.Config) (e *echo.Echo) {
	err := os.Mkdir(config.Server.UploadStorage, 0777)
	logrus.Fatal(err)

	server := echo.New()
	group := server.Group(dependencies.ConfigInstance().Server.APIPrefix,
		middleware.LoggerWithConfig(GetAccessConfig(config.Logs.Access)))

	group.POST("/sign_up", handlers.SignUp)
	group.POST("/sign_in", handlers.SignIn)
	group.POST("/upload", handlers.UploadImage)
	group.POST("/graphql", handlers.API)

	group.GET("/ping", handlers.Ping)
	group.GET("/ping_db", handlers.PingDb)

	AuthGroup := group.Group("")
	AuthGroup.Use(middleware.JWTWithConfig(handlers.GetJwtConfig()))
	AuthGroup.Use(handlers.LoadUser)
	AuthGroup.POST("/auth_ping", handlers.PingAuth)

	return server
}
