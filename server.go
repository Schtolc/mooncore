package main

import (
	"github.com/Schtolc/mooncore/handlers"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// InitServer inits echo server: sets access logs and handlers
func InitServer(config Config, db *gorm.DB) (e *echo.Echo) {
	server := echo.New()
	server.Use(middleware.LoggerWithConfig(GetAccessConfig(config.Logs.Access)))

	server.GET("/ping", handlers.Ping)
	server.GET("/ping_db", handlers.PingDb(db))

	return server
}
