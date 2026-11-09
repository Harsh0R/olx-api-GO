package httpx

// HTTP statuscode    Custom code
// 400                invalid_id
// 404                not_found
// 500                internal_error
// 400                malformed_json
// 422                validation_failed
// 401                unauthenticated
// 403                forbidden
// 409                conflict
// 429                rate_limited

import (
	"encoding/json"
	"net/http"
)

type CustomCode string

const (
	CODEInvalidId        CustomCode = "invalid_id"
	CODENotFound         CustomCode = "not_found"
	CODEInternalError    CustomCode = "internal_error"
	CODEMalformedJSON    CustomCode = "malformed_json"
	CODEValidationFailed CustomCode = "validation_failed"
	CODEUnauthenticated  CustomCode = "unauthenticated"
	CODEForbidden        CustomCode = "forbidden"
	CODEConflict         CustomCode = "conflict"
	CODERateLimited      CustomCode = "rate_limited"
)

type errorEnvelop struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Code    CustomCode `json:"code"`
	Message string     `json:"message"`
	Field   string     `json:"field,omitempty"`
}

func Error(resw http.ResponseWriter, status int, message string, code CustomCode) {
	resw.Header().Set("Content-Type", "application/json")
	resw.WriteHeader(status)
	_ = json.NewEncoder(resw).Encode(errorEnvelop{Error: errorPayload{
		Code:    code,
		Message: message,
	}})
}
func ValidationError(resw http.ResponseWriter, status int, message string, code CustomCode, field string) {
	resw.Header().Set("Content-Type", "application/json")
	resw.WriteHeader(status)
	_ = json.NewEncoder(resw).Encode(errorEnvelop{Error: errorPayload{
		Code:    code,
		Message: message,
		Field:   field,
	}})
}
