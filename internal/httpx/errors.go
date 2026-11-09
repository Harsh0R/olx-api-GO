package httpx

import (
	"encoding/json"
	"net/http"
)

func Error(resw http.ResponseWriter, status int, message string, code string) {
	resw.Header().Set("Content-Type", "application/json")
	resw.WriteHeader(status)
	_ = json.NewEncoder(resw).Encode(map[string]string{"message": message, "code": code})
}
