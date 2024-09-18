package util

import "go-service.codymj.io/internal/user"

type Services struct {
	//ValidatorService validate.Service
	UserService user.Service
}

type HttpResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
