package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// 加密密码
func HashPassword(password string) (string, error) {
	// 使用bcrypt.DefaultCost作为加密成本
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// 验证密码
func CheckPassword(password, hashedPassword string) bool {
	// 比较密码与哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}