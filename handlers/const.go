package handlers

var (
	needRegistration = &Resp{
		Code:    "400",
		Message: "You need to register",
	}
	invalidToken = &Resp{
		Code:    "400",
		Message: "Token is invalid",
	}
	internalError = &Resp{
		Code:    "400",
		Message: "Internal Error",
	}
	userAlreadyExists = &Resp{
		Code:    "400",
		Message: "User already exists in database",
	}
)
