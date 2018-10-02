package handlers

import (
	"../services"
	//"fmt"
	//"log"
	"net/http"
)

// Hello exposes an api for the hello service
type Hello struct {
	Service services.HelloService
}

// NewHello creates a new handler for hello
func NewHello(s services.HelloService) *Hello {
	return &Hello{s}
}

// Handler handles hello requests
func (h *Hello) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	//w.Header().Set("Access-Control-Allow-Headers", "Authorization, authorization, Content-Type")

	//session, err := store.Get(req, "session")

	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//data := session.Values["Auth"]
	//fmt.Print("Auth: ")
	//fmt.Println(data)
	//auth, ok := data.(string)
	//if ok != false {
	//	w.Header().Set("Authorization ", auth)
	//}

	switch req.Method {
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	case "GET":
		s := h.Service.SayHello()
		w.Write([]byte(s))
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
