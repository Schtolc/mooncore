package handlers

// Response model
type Response struct {
	Code int         `json:"code"`
	Body interface{} `json:"body"`
}
