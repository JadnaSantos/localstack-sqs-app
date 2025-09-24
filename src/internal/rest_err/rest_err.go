package resterr

import "net/http"

type RestErr struct {
	Message string `json:"message" example:"error trying to process request"`

	Err string `json:"error" example:"internal_server_error"`

	Code int `json:"code" example:"500"`

	Causes []Causes `json:"causes"`
}

type Causes struct {
	Field string `json:"field" example:"name"`

	Message string `json:"message" example:"name is required"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}
}
