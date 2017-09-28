package main

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/handlers"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// InitServer inits echo server: sets access logs and handlers
func InitServer(config *dependencies.Config, db *gorm.DB) (e *echo.Echo) {
	server := echo.New()
	server.Use(middleware.LoggerWithConfig(GetAccessConfig(config.Logs.Access)))

	server.POST("/sign_up", handlers.SignUp)
	server.POST("/sign_in", handlers.SignIn)

	server.GET("/ping", handlers.Ping)
	server.GET("/ping_db", handlers.PingDb)

	AuthGroup := server.Group("/")
	AuthGroup.Use(middleware.JWTWithConfig(handlers.GetJwtConfig()))
	AuthGroup.Use(handlers.CheckJwtToken)
	AuthGroup.POST("auth_ping", handlers.PingAuth)
	AuthGroup.GET("/graphql", handlers.API)


	return server
}
