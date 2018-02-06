package utils

import (
	"github.com/labstack/echo"
	"net/http"
)

// Response model
type Response struct {
	Data interface{} `json:"data"`
}

// SendResponse sends json response with given code and body
func SendResponse(c echo.Context, code int, body interface{}) error {
	return c.JSON(code, Response{body})
}

// InternalServerError send response with code = 500
func InternalServerError(c echo.Context) error {
	return SendResponse(c, http.StatusInternalServerError, "Internal server error")
}
