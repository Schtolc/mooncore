package handlers

var (
	needRegistration = &Resp{
		Code:    "400",
		Message: "You need to register",
	}
	invalidToken = &Resp{
		Code:    "Invalid Argument",
		Message: "Token is invalid",
	}
	internalError = &Resp{
		Code:    "InternalError",
		Message: "Internal Error",
	}
	userAlreadyExists = &Resp{
		Code:    "UserAlreadyExists",
		Message: "User already exists in database",
	}

)
