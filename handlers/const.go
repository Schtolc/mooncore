package handlers

var (
	// InternalError response
	InternalError = &Resp{
		Code:    "500",
		Message: "InternalError",
	}
	// NeedRegistration response
	NeedRegistration = &Resp{
		Code:    "400",
		Message: "You need to register",
	}
	// InvalidToken response
	InvalidToken = &Resp{
		Code:    "400",
		Message: "Token is invalid",
	}
)
