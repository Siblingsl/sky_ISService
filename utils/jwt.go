package utils

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"sky_ISService/config"
	"time"
)

// GenerateToken 生成 JWT Token
func GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub_id": userID,                                // 用户ID
		"role":   role,                                  // 用户角色
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // 72 小时过期
		"iat":    time.Now().Unix(),                     // 签发时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWTSecret.Secret))
}

// ParseToken 解析 JWT Token
// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// 打印 Token 内容，检查传入的 Token 是否有效
	log.Println("Token String:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 将密钥转换为 []byte 类型
		return []byte(config.GetConfig().JWTSecret.Secret), nil
	})
	if err != nil {
		// 打印解析错误
		log.Println("Error parsing token:", err)
		return nil, err
	}
	// 验证 token 是否有效并且转换为 jwt.MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 打印解析的 Claims
		log.Println("Parsed Claims:", claims)
		return claims, nil
	}

	// 如果解析失败，输出错误信息
	log.Println("Token is invalid or claims are of invalid type.")
	return nil, err
}
