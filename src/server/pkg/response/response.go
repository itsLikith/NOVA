package response

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

func SendSuccessResponse(status int, message string, data any) SuccessResponse {
	return SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func SendErrorResponse(status int, message string, error any) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Message: message,
		Error:   error,
	}
}
