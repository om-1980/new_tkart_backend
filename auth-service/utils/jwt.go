package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID 		int    `json:"user_id"`
	Identifier  string `json:"identifier"`
	Name   		string `json:"name"`
	Email  		string `json:"email"`
	Mobile 		string `json:"mobile"`
	Role   		string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, role, identifier, name, email, mobile string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Identifier: identifier,
		Name:   name,
		Email:  email,
		Mobile: mobile,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}


func ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
