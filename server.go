package main

import (
	"github.com/Schtolc/mooncore/api"
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

	server.GET("/ping", handlers.Ping)
	server.GET("/ping_db", handlers.PingDb(db))
	server.GET("/graphql", func(c echo.Context) error {
		result := api.ExecuteQuery(c.QueryParams().Get("query"), api.CreateSchema(db))
		return c.JSON(200, result)
	})

	return server
}
