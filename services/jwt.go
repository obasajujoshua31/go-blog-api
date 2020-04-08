package services

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-blog-api/config"
	"time"
)

var appConfig = config.LoadConfig()

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"iat":    time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(appConfig.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserIDFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(appConfig.JwtSecret), nil
	})
	if err != nil {
		return "", errors.New("unable to parse jwt")
	}

	if token == nil {
		return "", errors.New("token is nil")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["userID"].(string), nil
	} else {
		return "", errors.New("jwt invalid")
	}
}
