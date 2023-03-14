package model

import (
	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}
