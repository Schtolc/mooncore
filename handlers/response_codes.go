package handlers

var (
	// OK code
	OK = 200
	//Created code
	Created = 201
	// Accepted code
	Accepted = 202
	// NoContent code
	NoContent = 204

	// BadRequest code
	BadRequest = 400
	// Unauthorized code
	Unauthorized = 401
	// Forbidden code
	Forbidden = 403
	//NotFound code
	NotFound = 404
	//MethodNotAllowed code
	MethodNotAllowed = 405

	// InternalServerError code
	InternalServerError = 500
	// NotImplemented code
	NotImplemented = 501
	//BadGateway code
	BadGateway = 502
	// ServiceUnavailable code
	ServiceUnavailable = 503
)


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
	invalidParam = &Response{
		Code:   BadRequest,
		Body: "invalid request parameter",
	}

)