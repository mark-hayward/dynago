package services

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Set our secret.
// TODO: Use generated key from README
var mySigningKey = []byte("secret")

// Token defines a token for our application
type Token string

// TokenService provides a token
type TokenService interface {
	Get(u *User) (string, error)
}

type tokenService struct {
	UserService userService
}

// CustomClaims type holds the token claims
type CustomClaims struct {
	Admin bool  `json:"admin"`
	User  *User `json:"user"`
	jwt.StandardClaims
}

// NewTokenService creates a new UserService
func NewTokenService() TokenService {
	return &tokenService{}
}

// Get retrieves a token for a user
// TODO: Take login credentials and verify them against what's in database
func (s *tokenService) Get(u *User) (string, error) {

	// Try to log in the user
	user, err := s.UserService.Read(u.ID)
	if err != nil {
		return "", errors.New("Failed to retrieve user")
	}
	if user == nil {
		return "", errors.New("Failed to retrieve user")
	}

	claims := CustomClaims{
		Admin: true,
		User:  u,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with key
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", errors.New("Failed to sign token")
	}

	return tokenString, nil
}
