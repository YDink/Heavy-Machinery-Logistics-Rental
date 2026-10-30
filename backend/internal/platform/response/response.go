package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// ResponseMeta holds standard meta envelope fields conforming to AGENTS.md § 2.
type ResponseMeta struct {
	Version    string          `json:"version"`
	RequestID  string          `json:"requestId"`
	Timestamp  string          `json:"timestamp"`
	Pagination *PaginationMeta `json:"pagination"`
	Warnings   []string        `json:"warnings"`
}

type PaginationMeta struct {
	NextCursor *string `json:"nextCursor"`
	PrevCursor *string `json:"prevCursor"`
	HasMore    bool    `json:"hasMore"`
}

// ProblemDetails represents RFC 9457 error payload.
type ProblemDetails struct {
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	Status        int            `json:"status"`
	Detail        string         `json:"detail"`
	Code          string         `json:"code"`
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type SuccessEnvelope struct {
	Success bool         `json:"success"`
	Data    any          `json:"data"`
	Error   *string      `json:"error"`
	Meta    ResponseMeta `json:"meta"`
}

type ErrorEnvelope struct {
	Success bool            `json:"success"`
	Data    any             `json:"data"`
	Error   ProblemDetails  `json:"error"`
	Meta    ResponseMeta    `json:"meta"`
}

func getOrGenerateRequestID(r *http.Request) string {
	if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
		return reqID
	}
	return "req_" + uuid.NewString()[:16]
}

// JSON renders a unified success envelope.
func JSON(w http.ResponseWriter, r *http.Request, statusCode int, data any) {
	meta := ResponseMeta{
		Version:   "1.0",
		RequestID: getOrGenerateRequestID(r),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Warnings:  []string{},
	}

	envelope := SuccessEnvelope{
		Success: true,
		Data:    data,
		Error:   nil,
		Meta:    meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(envelope)
}

// Error renders an RFC 9457 Problem Details error wrapped in the unified envelope.
func Error(w http.ResponseWriter, r *http.Request, statusCode int, problem ProblemDetails) {
	problem.Status = statusCode
	if problem.Type == "" {
		problem.Type = "https://api.dinkmachinery.com/errors/" + problem.Code
	}

	meta := ResponseMeta{
		Version:   "1.0",
		RequestID: getOrGenerateRequestID(r),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Warnings:  []string{},
	}

	envelope := ErrorEnvelope{
		Success: false,
		Data:    nil,
		Error:   problem,
		Meta:    meta,
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(envelope)
}
