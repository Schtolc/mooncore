package handlers

var (
	InternalError = &Resp{
		Code:    "500",
		Message: "InternalError",
	}
	NeedRegistration = &Resp{
		Code:    "400",
		Message: "You need to register",
	}
	InvalidToken = &Resp{
		Code:    "400",
		Message: "Token is invalid",
	}
)
