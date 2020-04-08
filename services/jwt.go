package services

import (
	"github.com/dgrijalva/jwt-go"
	"go-blog-api/config"
	"time"
)

func GenerateJWT(userID string) (string, error) {
	appConfig := config.LoadConfig()
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
