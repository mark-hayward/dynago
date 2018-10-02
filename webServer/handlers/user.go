package handlers

import (
	"net/http"

	"../services"
)

// Users provides API access to the user service
type Users struct {
	Service services.UserService
}

// NewUsers creates a users.json handler
func NewUsers(s services.UserService) *Users {
	return &Users{s}
}

// Handler handles requests to /users.json
func (u *Users) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
	case "POST":
	case "PUT":
	case "DELETE":
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
