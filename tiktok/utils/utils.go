package utils

import "golang.org/x/crypto/bcrypt"

// EncryptPassword 对密码进行加密
func EncryptPassword(password string) (string, error) {
	cost := 5

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hashed), err
}

// VerifyPassword 验证密码
func VerifyPassword(pass, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	if err != nil {
		return false
	}
	return true
}
