package handlers

import (
	"github.com/labstack/echo"
	"net/http"
)

// Response model
type Response struct {
	Data interface{} `json:"data"`
}

func sendResponse(c echo.Context, code int, body interface{}) error {
	return c.JSON(code, Response{body})
}

func internalServerError(c echo.Context) error {
	return sendResponse(c, http.StatusInternalServerError, "Internal server error")
}
