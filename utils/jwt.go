package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 自定义声明结构体
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	// 设置过期时间，默认24小时
	expirationTime := time.Now().Add(24 * time.Hour)
	
	// 创建自定义声明
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-backend",
		},
	}
	
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名token，从环境变量获取密钥，没有则使用默认密钥（生产环境不建议）
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key" // 实际生产环境应使用更安全的密钥
	}
	
	// 生成签名字符串
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// 验证JWT令牌并返回声明
func ParseToken(tokenString string) (*Claims, error) {
	// 从环境变量获取密钥
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key"
	}
	
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// 验证claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}