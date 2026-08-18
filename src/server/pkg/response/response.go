package response

type SuccessResponse struct {
	Status int
	Message string
	Data any
}

type ErrorResponse struct {
	Status int
	Message string
	Error any
}

func SendSuccessResponse(status int, message string, data any) SuccessResponse {
	return SuccessResponse{
		Status: status,
		Message: message,
		Data: data,
	}
}

func SendErrorResponse(status int, message string, error any) ErrorResponse {
	return ErrorResponse{
		Status: status,
		Message: message,
		Error: error,
	}
}