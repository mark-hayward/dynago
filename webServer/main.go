package main

import (
	"log"
	"net/http"
)

func main() {
	certPath := "server.pem"
	keyPath := "server.key"
	api := NewAPI(certPath, keyPath)

	http.Handle("/hello", api.Hello)
	http.Handle("/table", api.Hello)
	http.Handle("/tables", api.Hello)
	http.Handle("/tokens", api.Tokens)
	http.Handle("/", http.FileServer(http.Dir("../web")))
	http.Handle("/process", AddMiddleware(http.HandlerFunc(Collector), api.Authenticate))

	http.Handle("/encryptTest", AddMiddleware(api.Test, api.Authenticate))

	http.Handle("/users.json", AddMiddleware(api.Users, api.Authenticate))

	err := http.ListenAndServeTLS(":8001", certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// AddMiddleware adds middleware to a Handler
func AddMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}
