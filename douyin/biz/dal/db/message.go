package db

import (
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"time"
)

type Messages struct {
	ID         int64     `json:"id"`
	ToUserId   int64     `json:"to_user_id"`
	FromUserId int64     `json:"from_user_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Messages) TableName() string {
	return constants.MessageTableName
}

// AddNewMessage 增加新消息
func AddNewMessage(msg *Messages) (bool, error) {
	exist, err := QueryUserById(msg.FromUserId)
	if exist == nil || err != nil {
		return false, errno.UserIsNotExistErr
	}
	exist, err = QueryUserById(msg.ToUserId)
	if exist == nil || err != nil {
		return false, errno.UserIsNotExistErr
	}
	err = dbConn.Create(msg).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetLatestMessageByIdPair 查询usr1和usr2的最新消息
func GetLatestMessageByIdPair(uid1, uid2 int64) (*Messages, error) {
	exist, err := QueryUserById(uid1)
	if exist == nil || err != nil {
		return nil, errno.UserIsNotExistErr
	}
	exist, err = QueryUserById(uid2)
	if exist == nil || err != nil {
		return nil, errno.UserIsNotExistErr
	}

	message := Messages{}
	err = dbConn.Where("to_user_id = ? AND from_user_id = ?", uid1, uid2).Or(
		"to_user_id = ? AND from_user_id = ?", uid1, uid2).Last(&message).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}

	return &message, err
}
