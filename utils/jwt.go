package utils

import (
	"github.com/dgrijalva/jwt-go"
	"sky_ISService/config"
	"time"
)

// GenerateToken 生成 JWT Token
func GenerateToken(userID string, role string) (string, error) {

	claims := jwt.MapClaims{
		"sub_id": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // 72 小时过期
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWTSecret.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetConfig().JWTSecret.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
