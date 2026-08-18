package services

import (
	"github.com/nova/pkg/response"
)

func HealthStatus() response.SuccessResponse {
	return response.SuccessResponse{
		Status:  200,
		Message: "Auth service is UP",
		Data:    nil,
	}
}