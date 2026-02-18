package authorization

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("secret-key") //секретный ключ, которым подписывается токен

type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
} //Это структура данных, которая будет храниться внутри JWT (payload). Она содержит: UserID — кастомное поле (ID пользователя) jwt.RegisteredClaims — стандартные JWT-поля: exp — срок действия (expiration) iat — время создания (issued at) iss — издатель

// принимает userID/возвращает строку (JWT) и ошибку, если что-то пошло не так
func GenerateToken(user_id int) (string, error) {
	claims := Claims{
		UserId: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //Токен будет действовать 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "note_service",
			Subject:   fmt.Sprint(user_id),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //Создание токена
	return token.SignedString(JwtKey)
}
