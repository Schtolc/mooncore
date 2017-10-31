package handlers

import (
	"github.com/labstack/echo"
	"net/http"
)

// Response model
type Response struct {
	Code int         `json:"code"`
	Body interface{} `json:"body"`
}

func sendResponse(c echo.Context, code int, body interface{}) error {
	return c.JSON(http.StatusOK, Response{code, body})
}

func internalServerError(c echo.Context) error {
	return sendResponse(c, http.StatusInternalServerError, "Internal server error")
}
