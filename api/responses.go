package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSONResponse(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func ErrResponse(w http.ResponseWriter,
	_ *http.Request,
	errResponse ErrorResponse,
) {
	if errResponse.Code == 0 {
		errResponse.Code = http.StatusInternalServerError
	}
	if errResponse.Error == "" {
		errResponse.Error = http.StatusText(http.StatusInternalServerError)
	}
	JSONResponse(w, errResponse.Code, errResponse)
}
