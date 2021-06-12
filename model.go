package main 
import (
	"github.com/dgrijalva/jwt-go"
)
type SignupJSON struct {
	Name        string  `json:"name"`
	Rollno      int64   `json:"rollno"`
	Password    string    `json:"password"`
}
type LoginJSON struct {
	Rollno      int64   `json:"rollno"`
	Password    string    `json:"password"`
}

type CustomClaims struct {
	Rollno      int64   `json:"rollno"`
	jwt.StandardClaims
}

type LoginToken struct {
	Name string `json:"name"`
	Token string `json:"token"`
}