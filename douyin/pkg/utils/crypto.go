package utils

import "golang.org/x/crypto/bcrypt"

// CryptPassword 对密码进行加密
func CryptPassword(password string) (string, error) {
	cost := 5

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hashed), err
}

// VerifyPassword 验证密码
func VerifyPassword(pwd, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}
