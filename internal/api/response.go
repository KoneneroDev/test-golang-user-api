package api

import "net/http"

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OkResponse() Response {
	return Response{
		Status: http.StatusOK,
	}
}

func ErrorStatus(msg string) Response {
	return Response{
		Status: http.StatusBadRequest,
		Error:  msg,
	}
}
