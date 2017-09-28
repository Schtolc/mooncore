package main

import (
	"github.com/Schtolc/mooncore/handlers"
	"github.com/Schtolc/mooncore/utils"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// InitServer inits echo server: sets access logs and handlers
func InitServer(config utils.Config, db *gorm.DB) (e *echo.Echo) {
	server := echo.New()
	server.Use(middleware.LoggerWithConfig(GetAccessConfig(config.Logs.Access)))

	h := handlers.Init(db)

	server.POST("/sign_up", h.SignUp)
	server.POST("/sign_in", h.SignIn)

	server.GET("/ping", h.Ping)
	server.GET("/ping_db", h.PingDb)

	AuthGroup := server.Group("/")
	AuthGroup.Use(middleware.JWTWithConfig(handlers.GetJwtConfig()))
	AuthGroup.Use(h.CheckJwtToken)
	AuthGroup.POST("auth_ping", h.PingAuth)

	return server
}
