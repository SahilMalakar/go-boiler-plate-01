package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

// custome field for error code
const (
	CodeInvalidID        Code = "invalid id"
	CodeNotFound         Code = "not found"
	CodeInternalError    Code = "internal error"
	CodeMalformedJSON    Code = "malformed json"
	CodeValidationFailed Code = "validation failed"
	CodeUnauthenticated  Code = "unauthenticated"
	CodeForbidden        Code = "forbidden"
	CodeConflict         Code = "conflict"
	CodeRateLimited      Code = "rate limited"
)

type errorPayload struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

func Error(w http.ResponseWriter, status int, msg string, code Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(
		errorEnvelope{
			Error: errorPayload{
				Code:    code,
				Message: msg,
			},
		},
	)
}
