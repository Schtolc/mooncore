package handlers

// Response model
type Response struct {
	Code int         `json:"code"`
	Body interface{} `json:"body"`
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