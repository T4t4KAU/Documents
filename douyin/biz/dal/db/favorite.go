package db

import (
	"douyin/biz/mw/redis"
	"douyin/pkg/constants"
	"gorm.io/gorm"
	"time"
)

var rdFavorite redis.Favorite

type Favorites struct {
	ID        int64          `json:"id"`
	UserId    int64          `json:"user_id"`
	VideoId   int64          `json:"video_id"`
	CreatedAt time.Time      `json:"create_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

func (Favorites) TableName() string {
	return constants.FavoritesTableName
}

// AddNewFavorite 添加点赞信息
func AddNewFavorite(favorite *Favorites) (bool, error) {
	if err := dbConn.Create(favorite).Error; err != nil {
		return false, err
	}

	// 将点赞数据添加到redis
	if rdFavorite.CheckLiked(favorite.VideoId) {
		rdFavorite.AddLiked(favorite.UserId, favorite.VideoId)
	}

	if rdFavorite.CheckLike(favorite.UserId) {
		rdFavorite.AddLike(favorite.UserId, favorite.VideoId)
	}

	return true, nil
}

// DeleteFavorite 删除评论
func DeleteFavorite(favorite *Favorites) (bool, error) {
	if err := dbConn.Where("video_id = ? AND user_id = ?",
		favorite.VideoId, favorite.UserId).Delete(favorite).Error; err != nil {
		return false, err
	}

	if rdFavorite.CheckLiked(favorite.VideoId) {
		rdFavorite.DelLiked(favorite.UserId, favorite.VideoId)
	}

	if rdFavorite.CheckLike(favorite.UserId) {
		rdFavorite.DelLiked(favorite.UserId, favorite.VideoId)
	}

	return true, nil
}

// QueryFavoriteExist 查询点赞是否存在
func QueryFavoriteExist(userId, videoId int64) (bool, error) {
	if rdFavorite.CheckLike(videoId) {
		return rdFavorite.ExistLiked(userId, videoId), nil
	}
	if rdFavorite.CheckLiked(userId) {
		return rdFavorite.ExistLike(userId, videoId), nil
	}

	var sum int64
	err := dbConn.Model(&Favorites{}).Where(
		"video_id ? AND user_id = ?", videoId, userId).Count(&sum).Error
	if err != nil {
		return false, err
	}
	if sum == 0 {
		return false, nil
	}
	return true, nil
}

// QueryTotalFavoritedByAuthorId 查询点赞总数
func QueryTotalFavoritedByAuthorId(authorId int64) (int64, error) {
	var sum int64
	err := dbConn.Table(constants.FavoritesTableName).Joins(
		"JOIN videos ON favorites.video_id = videos.id").Where(
		"videos.author_id = ?", authorId).Count(&sum).Error
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// 获取指定用户ID的点赞列表
func getFavoriteIdList(userId int64) ([]int64, error) {
	var favoriteActions []Favorites

	err := dbConn.Where("user_id = ?", userId).Find(&favoriteActions).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range favoriteActions {
		result = append(result, v.VideoId)
	}

	return result, nil
}

// GetFavoriteCountByUserId 获取用户的点赞数
func GetFavoriteCountByUserId(userId int64) (int64, error) {
	if rdFavorite.CheckLike(userId) {
		return rdFavorite.CountLike(userId)
	}

	videoIds, err := getFavoriteIdList(userId)
	if err != nil {
		return 0, err
	}

	go func(uid int64, vids []int64) {
		for _, vid := range vids {
			rdFavorite.AddLiked(uid, vid)
		}
	}(userId, videoIds)

	return int64(len(videoIds)), nil
}

// GetFavoriteIdList 获取用户点赞列表
func GetFavoriteIdList(userId int64) ([]int64, error) {
	if rdFavorite.CheckLiked(userId) {
		return rdFavorite.GetLiked(userId), nil
	}
	return getFavoriteIdList(userId)
}

// GetFavoriteCount 获取视频点赞数
func GetFavoriteCount(videoId int64) (int64, error) {
	if rdFavorite.CheckLiked(videoId) {
		return rdFavorite.CountLiked(videoId)
	}

	favoriteIds, err := getFavoriteIdList(videoId)
	if err != nil {
		return 0, err
	}

	// 异步更新redis
	go func(users []int64, vid int64) {
		for _, u := range users {
			rdFavorite.AddLiked(u, videoId)
		}
	}(favoriteIds, videoId)

	return int64(len(favoriteIds)), nil
}
