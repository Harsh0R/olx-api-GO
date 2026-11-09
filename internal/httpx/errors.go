package httpx

import (
	"encoding/json"
	"net/http"
)

type errorEnvelop struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(resw http.ResponseWriter, status int, message string, code string) {
	resw.Header().Set("Content-Type", "application/json")
	resw.WriteHeader(status)
	_ = json.NewEncoder(resw).Encode(errorEnvelop{Error: errorPayload{
		Code:    code,
		Message: message,
	}})
}
