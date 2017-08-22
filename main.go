package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Resp struct {
	Code    string `json:"code" xml:"code"`
	Message string `json:"message" xml:"message"`
}

func main() {
	e := echo.New()
	e.Use(middleware.RequestID())

	db := ConnectDB()
	defer db.Close()

	e.GET("/ping", Ping)
	e.GET("/ping_db", func(c echo.Context) error {
		m := &Metric{
			Path: c.Path(),
			Time: time.Now(),
		}
		db.Create(m)
		return c.JSON(http.StatusOK, &Resp{
			Code:    "200",
			Message: m.Id,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, &Resp{
		Code:    "200",
		Message: "Hello, World!",
	})
}
