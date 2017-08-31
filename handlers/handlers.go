package handlers

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"mooncore/models"
	"net/http"
)

type Resp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, &Resp{
		Code:    "200",
		Message: "Hello, World!",
	})
}

func PingDb(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		m := &models.Metric{
			Path: c.Path(),
			Time: gorm.NowFunc(),
		}
		db.Create(m)
		return c.JSON(http.StatusOK, &Resp{
			Code:    "200",
			Message: m.Id,
		})
	}
}
