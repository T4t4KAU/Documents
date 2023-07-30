package db

import (
	"douyin/pkg/constants"
	"douyin/pkg/errno"
)

type User struct {
	ID              int64  `json:"id"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
}

func (User) TableName() string {
	return constants.UserTableName
}

// CreateUser 创建用户
func CreateUser(user *User) (int64, error) {
	err := dbConn.Create(user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, err
}

// QueryUserByName 通过用户名查询用户
func QueryUserByName(uname string) (*User, error) {
	var user User
	err := dbConn.Where("user_name = ?", uname).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// QueryUserById 通过用户ID查询用户
func QueryUserById(userId int64) (*User, error) {
	var user User
	err := dbConn.Where("id = ?", userId).Find(&user).Error
	if err != nil {
		return nil, err
	}

	if user == (User{}) {
		err := errno.UserIsNotExistErr
		return nil, err
	}
	return &user, nil
}

// VerifyUser 验证用户的用户名和密码
func VerifyUser(username, password string) (int64, error) {
	var user User
	err := dbConn.Where("user_name = ? AND password = ?",
		username, password).Find(&user).Error
	if err != nil {
		return 0, err
	}

	if user.ID == 0 {
		err := errno.PasswordIsNotVerified
		return user.ID, err
	}
	return user.ID, nil
}

// CheckUserExistById 通过用户ID检查用户是否存在
func CheckUserExistById(userId int64) (bool, error) {
	var user User
	err := dbConn.Where("id = ?", userId).Find(&user).Error
	if err != nil {
		return false, err
	}

	if user == (User{}) {
		return false, nil
	}
	return true, nil
}
