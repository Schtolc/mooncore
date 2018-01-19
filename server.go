package main

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/graphql"
	"github.com/Schtolc/mooncore/rest"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"github.com/Schtolc/mooncore/utils"
)

// InitServer inits echo server: sets access logs and graphql
func InitServer(config *dependencies.Config) (e *echo.Echo) {
	_ = os.Mkdir(config.Server.UploadStorage, 0777)

	server := echo.New()
	group := server.Group(dependencies.ConfigInstance().Server.APIPrefix,
		middleware.LoggerWithConfig(GetAccessConfig(config.Logs.Access)),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
				c.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, "content-type")
				c.Response().Header().Set(echo.HeaderAccessControlAllowMethods, "DELETE, GET, OPTIONS, PATCH, POST, PUT")
				c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
				return next(c)
			}
		},
		middleware.JWTWithConfig(utils.GetJwtConfig()),
		rest.LoadUser)

	group.POST("/upload", rest.UploadImage)
	group.POST("/graphql", graphql.API)
	group.OPTIONS("/graphql", rest.Headers)

	group.GET("/ping", rest.Ping)
	group.GET("/ping_db", rest.PingDb)

	group.POST("/auth_ping", rest.PingAuth)

	return server
}
