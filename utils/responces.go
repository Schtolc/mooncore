package utils

import (
	"github.com/labstack/echo"
	"net/http"
)

// Response model
type Response struct {
	Data interface{} `json:"data"`
}

func SendResponse(c echo.Context, code int, body interface{}) error {
	return c.JSON(code, Response{body})
}

func InternalServerError(c echo.Context) error {
	return SendResponse(c, http.StatusInternalServerError, "Internal server error")
}
