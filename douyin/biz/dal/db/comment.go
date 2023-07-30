package db

import (
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID          int64          `json:"id"`
	UserId      int64          `json:"user_id"`
	VideoId     int64          `json:"video_id"`
	CommentText string         `json:"comment_text"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Comment) TableName() string {
	return constants.CommentTableName
}

// AddNewComment 添加新评论
func AddNewComment(comment *Comment) error {
	if ok, _ := CheckUserExistById(comment.UserId); ok {
		return errno.UserIsNotExistErr
	}

	if ok, _ := CheckVideoExistById(comment.VideoId); ok {
		return errno.VideoIsNotExistErr
	}

	err := dbConn.Create(comment).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteCommentById 通过评论ID删除评论
func DeleteCommentById(cid int64) error {
	comment := Comment{}
	err := dbConn.Where("id = ?", cid).Find(&comment).Error
	if err != nil {
		return errno.CommentIsNotExistErr
	}
	return dbConn.Delete(&comment).Error
}

// CheckCommentExist 检查评论是否存在
func CheckCommentExist(cid int64) (bool, error) {
	comment := Comment{}
	err := dbConn.Where("id = ?", cid).Find(&comment).Error
	if err != nil {
		return false, err
	}
	if comment.ID == 0 {
		return false, nil
	}
	return true, nil
}

// GetCommentListByVideoId 通过视频ID获取视频列表
func GetCommentListByVideoId(vid int64) ([]*Comment, error) {
	var comments []*Comment
	if ok, _ := CheckVideoExistById(vid); !ok {
		return comments, errno.VideoIsNotExistErr
	}
	err := dbConn.Table(constants.CommentTableName).Where("video_id = ?", vid).Error
	if err != nil {
		return comments, err
	}
	return comments, nil
}

// GetCommentCountByVideoId 通过视频ID获取视频评论数目
func GetCommentCountByVideoId(vid int64) (int64, error) {
	var sum int64

	err := dbConn.Model(&Comment{}).Where("video_id = ?", vid).Count(&sum).Error
	if err != nil {
		return sum, err
	}
	return sum, nil
}
