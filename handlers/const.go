package handlers

var (
	// needRegistration  response
	needRegistration = &Resp{
		Code:    "400",
		Message: "You need to register",
	}
	//InvalidToken response
	invalidToken = &Resp{
		Code:    "Invalid Argument",
		Message: "Token is invalid",
	}
	// InternalError response
	internalError = &Resp{
		Code:    "InternalError",
		Message: "Internal Error",
	}

)
