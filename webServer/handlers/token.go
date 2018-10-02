package handlers

import (
	"../services"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	//"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
	"os"
)

//var secret, _ = hex.DecodeString("VERY-BIG-SECRET") //Import something via an environment variable

//var store = sessions.NewCookieStore(secret)

// Tokens exposes an API to the tokens service
type Tokens struct {
	Service services.TokenService
}

// NewTokens creates new handler for tokens
func NewTokens(s services.TokenService) *Tokens {
	return &Tokens{s}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func hash(s string) (h [16]byte) {
	h = md5.Sum([]byte(s))
	return
}

// Handler will return tokens
func (t *Tokens) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	//w.Header().Set("Access-Control-Allow-Headers", "Authorization, authorization, Content-Type")
	switch req.Method {
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	case "POST":
		// TODO: Open Config File
		// TODO: Take in login information from the req body
		// TODO: Take the Password from the request and hash it using md5
		// TODO: Compare the hashed Password in the header with the one in the config file
		// TODO: If passed in Username = json Username AND passed in Password (hashed) = json Password hash
		// TODO: Then create a jsonUser :
		//session, err := store.Get(req, "session")
		//check(err)

		//session.Options = &sessions.Options{
		//	Path:     "/",
		//	MaxAge:   86400, //24 hours
		//	HttpOnly: false,
		//}

		file, err := os.Open("config/users.json")
		check(err)

		//Use if there are > 1 users
		//type Users struct {
		//	Users  []jsonUser 'json:"users"'
		//}
		type clientCredentials struct {
			Username string `json:"Username"`
			Password string `json:"Password"`
		}

		type jsonUser struct {
			Username  string `json:"Username"`
			Password  string `json:"Password"`
			UserID    int    `json:"userID"`
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
		}

		byteValue, _ := ioutil.ReadAll(file)

		var user jsonUser

		json.Unmarshal(byteValue, &user)

		credentialValue, _ := ioutil.ReadAll(req.Body)

		var creds clientCredentials

		json.Unmarshal(credentialValue, &creds)
		//for i := 0; i < len(users.Users); i++ {
		//}
		decodedPass, _ := hex.DecodeString(user.Password)
		x := hash(creds.Password)
		hashedPass := x[:]

		if user.Username == creds.Username && bytes.Equal(hashedPass, decodedPass) {
			currentUser := &services.User{
				ID:        user.UserID,
				FirstName: user.Firstname,
				LastName:  user.Lastname,
				Roles:     []string{services.AdministratorRole},
			}
			if err != nil {
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			}
			token, err := t.Service.Get(currentUser)
			check(err)
			//jsonToken := "{\"AuthToken\": \"" + ("Bearer " + token) + "\"}"
			//fmt.Println(jsonToken)

			type JSONFormat struct {
				Token string `json:"token"`
			}
			group := &JSONFormat{
				Token: "Bearer " + token,
			}
			jsonToken, err := json.Marshal(group)
			check(err)
			//fmt.Println(group)
			//fmt.Println(jsonToken)
			//fmt.Println(string(jsonToken))
			w.Write([]byte(string(jsonToken)))

		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
		//session.Save(req, w)
		//If Username and Password returned an unauthorised message like :
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
