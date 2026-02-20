package main

import (
	"fmt"
	"log"
	"note_service/internal/authorization"

	"github.com/golang-jwt/jwt/v5"
)

// Данный код просто для проверки токена. Вписал сюда, когда не получалось авторизоваться в постмене
func main() {
	// Вставьте сюда ваш токен
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJpc3MiOiJub3RlX3NlcnZpY2UiLCJzdWIiOiIxIiwiZXhwIjoxNzcxNzE1NDQxLCJpYXQiOjE3NzE2MjkwNDF9.uzxd7r-9DJv3cBTTQreK0bz69bzgabboqTxDpfMdDtk"

	// Разбор токена с кастомными claims
	token, err := jwt.ParseWithClaims(tokenString, &authorization.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return authorization.JwtKey, nil
	})

	if err != nil {
		log.Fatalf("Ошибка при разборе токена: %v", err)
	}

	if !token.Valid {
		log.Fatalf("Токен недействителен")
	}

	claims, ok := token.Claims.(*authorization.Claims)
	if !ok || claims == nil {
		log.Fatalf("Не удалось привести claims к типу Claims")
	}

	fmt.Println("Токен действителен!")
	fmt.Println("UserID:", claims.UserId)
	fmt.Println("IssuedAt:", claims.IssuedAt)
	fmt.Println("ExpiresAt:", claims.ExpiresAt)
	fmt.Println("Issuer:", claims.Issuer)
}
