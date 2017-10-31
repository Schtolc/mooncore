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

var (
	needRegistration = &Response{
		Code:  Forbidden,
		Body: "You need to register",
	}
	invalidToken = &Response{
		Code: Forbidden,
		Body: "Token is invalid",
	}
	internalError = &Response{
		Code:  InternalServerError,
		Body: "Internal Error",
	}
	userAlreadyExists = &Response{
		Code: BadRequest,
		Body: "User already exists in database",
	}
	requireFields = &Response{
		Code:   BadRequest,
		Body: "require parameters for method",
	}
)