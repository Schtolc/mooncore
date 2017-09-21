package server

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/handlers"
	"github.com/Schtolc/mooncore/logger"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Init echo server: sets access logs and handlers
func Init(config config.Config, db *gorm.DB) (e *echo.Echo) {
	server := echo.New()
	server.Use(middleware.LoggerWithConfig(logger.GetAccessConfig(config.Logs.Access)))

	server.GET("/ping", handlers.Ping)
	server.GET("/ping_db", handlers.PingDb(db))

	return server
}
