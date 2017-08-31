package handlers

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"mooncore/models"
	"net/http"
)

type Resp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func EchoPing(c echo.Context) error {
	content, err := ioutil.ReadAll(c.Request().Body)
	check_err(err)
	return c.JSON(http.StatusOK, &Resp{
		Code:    "200",
		Message: string(content),
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
func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
