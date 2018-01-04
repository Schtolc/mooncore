package handlers

import (
	"github.com/labstack/echo"
	"net/http"
)


func sendResponse(c echo.Context, code int, body interface{}) error {
	switch body.(type) {
		case string:
			return c.String(code, body.(string));
		default:
			return c.JSON(code, body)
	}
}

func internalServerError(c echo.Context) error {
	return sendResponse(c, http.StatusInternalServerError, "Internal server error")
}
