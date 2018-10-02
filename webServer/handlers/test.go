package handlers

import (
	"../services"
	"net/http"
)

// Test exposes an api for the test service
type Test struct {
	Service services.TestService
}

// NewTest creates a new handler for test
func NewTest(s services.TestService) *Test {
	return &Test{s}
}

// Handler handles test requests
func (h *Test) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, authorization, Content-Type")
	switch req.Method {
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	case "GET":
		s := h.Service.TestResults()
		w.Write([]byte(s))
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
