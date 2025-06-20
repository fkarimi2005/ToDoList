package utils

import (
	"ToDoList/internal/configs"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	UserRole string `json:"user_role"`
	jwt.StandardClaims
}

func GenerateToken(userID int, username, UserRole string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		UserRole: UserRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().
				Add(time.Duration(configs.AppSettings.AuthParams.JwtTtlMinutes) * time.Minute).
				Unix(),
			Issuer:   configs.AppSettings.AppParams.ServerName,
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
