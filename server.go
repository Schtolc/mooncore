package main

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

// InitServer inits echo server: sets access logs and handlers
func InitServer(config *dependencies.Config) (e *echo.Echo) {
	_ = os.Mkdir(config.Server.UploadStorage, 0777)

	server := echo.New()
	group := server.Group(dependencies.ConfigInstance().Server.APIPrefix,
		middleware.LoggerWithConfig(GetAccessConfig(config.Logs.Access)),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
				return next(c)
			}
		})

	group.POST("/sign_up", handlers.SignUp)
	group.POST("/sign_in", handlers.SignIn)
	group.POST("/upload", handlers.UploadImage)
	group.POST("/graphql", handlers.API)
	group.OPTIONS("/graphql", handlers.Headers)

	group.GET("/ping", handlers.Ping)
	group.GET("/ping_db", handlers.PingDb)

	AuthGroup := group.Group("")
	AuthGroup.Use(middleware.JWTWithConfig(handlers.GetJwtConfig()))
	AuthGroup.Use(handlers.LoadUser)
	AuthGroup.POST("/auth_ping", handlers.PingAuth)

	return server
}
