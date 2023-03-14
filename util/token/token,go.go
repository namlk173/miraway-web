package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-mirayway/model"
	"time"
)

func CreateAccessToken(user *model.User, secret string, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &model.JwtCustomClaims{
		ID:       user.ID.Hex(),
		Username: user.UserName,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func CreateRefreshToken(user *model.User, secret string, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &model.JwtCustomRefreshClaims{
		ID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims["id"].(string), nil
}
