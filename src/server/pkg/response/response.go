package response

type Response struct {
	Message string      `json:"message"`
	Data    any `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewResponse(message string, data any, err string) *Response {
	return &Response{
		Message: message,
		Data:    data,
		Error:   err,
	}
}