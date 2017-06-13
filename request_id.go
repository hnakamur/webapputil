package webapputil

import (
	"context"
	"net/http"
)

// RequestIDFunc is a function type for generating a request ID.
type RequestIDFunc func(req *http.Request) string

type requestIDKeyType struct{}

var requestIDKey = requestIDKeyType{}

// ReqestID returns the request ID which was set by RequestIDMiddleware.
func RequestID(r *http.Request) string {
	return r.Context().Value(requestIDKey).(string)
}

// RequestIDMiddleware is a http middleware to generate a request ID
// and set it to the request context.
func RequestIDMiddleware(next http.Handler, fn RequestIDFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := fn(r)
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
