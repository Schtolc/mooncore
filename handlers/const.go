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
)
