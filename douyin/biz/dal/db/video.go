package db

import (
	"douyin/pkg/constants"
	"time"
)

type Video struct {
	ID          int64     `json:"id"`
	AuthorID    int64     `json:"author_id"`
	PlayURL     string    `json:"play_url"`
	CoverURL    string    `json:"cover_url"`
	PublishTime time.Time `json:"publish_time"`
	Title       string    `json:"title"`
}

func (Video) TableName() string {
	return constants.VideosTableName
}

// CreateVideo 创建video
func CreateVideo(video *Video) (int64, error) {
	err := dbConn.Create(video).Error
	if err != nil {
		return 0, err
	}
	return video.ID, err
}

// GetVideosByLastTime 通过指定时间获取视频列表
func GetVideosByLastTime(lastTime time.Time) ([]*Video, error) {
	videos := make([]*Video, constants.VideoFeedCount)
	err := dbConn.Where("publish_time < ?", lastTime).Error
	if err != nil {
		return videos, err
	}
	return videos, nil
}

// GetVideoByUserId 通过用户ID获取视频
func GetVideoByUserId(userId int64) ([]*Video, error) {
	var videos []*Video
	err := dbConn.Where("author_id = ?", userId).Find(&videos).Error
	if err != nil {
		return videos, err
	}
	return videos, nil
}

// GetVideosByVideoIdList 通过视频ID列表获取视频列表
func GetVideosByVideoIdList(videoIds []int64) ([]*Video, error) {
	var videos []*Video

	for _, vid := range videoIds {
		var video *Video
		err := dbConn.Where("id = ?", vid).Find(&video).Error
		if err != nil {
			return videos, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// GetPublishCountById 通过用户ID获取用户发布视频数
func GetPublishCountById(userId int64) (int64, error) {
	var count int64

	err := dbConn.Model(&Video{}).Where("author_id = ?", userId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CheckVideoExistById 通过视频ID检查视频是否存在
func CheckVideoExistById(videoId int64) (bool, error) {
	var video Video

	err := dbConn.Where("id = ?", videoId).Find(&video).Error
	if err != nil {
		return false, err
	}
	if video == (Video{}) {
		return false, nil
	}
	return true, nil
}
