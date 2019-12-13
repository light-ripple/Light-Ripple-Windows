package common

// ResponseBase is the data that is always returned with an API request.
type ResponseBase struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// GetCode retrieves the response code.
func (r ResponseBase) GetCode() int {
	return r.Code
}

// SetCode changes the response code.
func (r *ResponseBase) SetCode(i int) {
	r.Code = i
}

// GetMessage retrieves the response message.
func (r ResponseBase) GetMessage() string {
	return r.Message
}

// CodeMessager is something that has the Code() and Message() methods.
type CodeMessager interface {
	GetMessage() string
	GetCode() int
}

// SimpleResponse returns the most basic response.
func SimpleResponse(code int, message string) CodeMessager {
	return ResponseBase{
		Code:    code,
		Message: message,
	}
}
