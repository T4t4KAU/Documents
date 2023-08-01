package db

import (
	"context"
	"douyin/cmd/relation/dal/redis"
	"douyin/pkg/constants"
	"gorm.io/gorm"
	"time"
)

var rdFollows redis.Follows

// Relation  用户关系
type Relation struct {
	ID         int64          `json:"id"`
	UserId     int64          `json:"user_id"`
	FollowerId int64          `json:"follower_id"`
	CreatedAt  time.Time      `json:"create_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

func (Relation) TableName() string {
	return constants.RelationTableName
}

// AddNewRelation 添加用户关系
func AddNewRelation(relation *Relation) (bool, error) {
	err := dbConn.Create(relation).Error
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	if rdFollows.CheckFollow(ctx, relation.FollowerId) {
		rdFollows.AddFollow(ctx, relation.UserId, relation.FollowerId)
	}

	if rdFollows.CheckFollower(ctx, relation.UserId) {
		rdFollows.AddFollower(ctx, relation.UserId, relation.FollowerId)
	}

	return true, nil
}

// DeleteRelation 删除用户关系
func DeleteRelation(relation *Relation) (bool, error) {
	err := dbConn.Where("user_id = ? AND follower_id = ?",
		relation.UserId, relation.FollowerId).Delete(relation).Error
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	if rdFollows.CheckFollow(ctx, relation.FollowerId) {
		rdFollows.DelFollower(ctx, relation.UserId, relation.FollowerId)
	}

	if rdFollows.CheckFollower(ctx, relation.UserId) {
		rdFollows.DelFollower(ctx, relation.UserId, relation.FollowerId)
	}

	return true, nil
}

// CheckRelationExist 查询关系是否存在
func CheckRelationExist(ctx context.Context, userId, followerId int64) (bool, error) {
	var r Relation

	if rdFollows.CheckFollow(ctx, followerId) {
		return rdFollows.ExistFollow(ctx, userId, followerId), nil
	}
	if rdFollows.CheckFollower(ctx, userId) {
		return rdFollows.ExistFollower(ctx, userId, followerId), nil
	}

	err := dbConn.Where("user_id = ? AND follower_id = ?", userId, followerId).Find(&r).Error
	if err != nil {
		return false, err
	}
	if r.ID == 0 {
		return false, nil
	}

	return true, nil
}

// GetFollowCount 查询用户关注数量
func GetFollowCount(ctx context.Context, followerId int64) (int64, error) {
	if rdFollows.CheckFollow(ctx, followerId) {
		return rdFollows.CountFollow(ctx, followerId)
	}

	// 缓存未命中 访问数据库
	followings, err := getFollowIdList(followerId)
	if err != nil {
		return 0, err
	}

	// 更新redis
	go addFollowRelationToRedis(followerId, followings)
	return int64(len(followings)), nil
}

// GetFollowerCount 查询用户粉丝数量
func GetFollowerCount(ctx context.Context, userId int64) (int64, error) {
	if rdFollows.CheckFollower(ctx, userId) {
		return rdFollows.CountFollower(ctx, userId)
	}

	followers, err := getFollowerIdList(userId)
	if err != nil {
		return 0, err
	}

	go addFollowerRelationToRedis(userId, followers)
	return int64(len(followers)), nil
}

// GetFollowIdList 获取用户关注列表
func GetFollowIdList(ctx context.Context, followerId int64) ([]int64, error) {
	if rdFollows.CheckFollow(ctx, followerId) {
		return rdFollows.GetFollow(ctx, followerId), nil
	}

	return getFollowIdList(followerId)
}

// GetFollowerIdList 获取用户粉丝列表
func GetFollowerIdList(ctx context.Context, followerId int64) ([]int64, error) {
	if rdFollows.CheckFollower(ctx, followerId) {
		return rdFollows.GetFollow(ctx, followerId), nil
	}

	return getFollowIdList(followerId)
}

func getFollowIdList(followerId int64) ([]int64, error) {
	var relations []Relation
	err := dbConn.Where("follower_id = ?", followerId).Find(&relations).Error
	if err != nil {
		return nil, err
	}

	var result []int64
	for _, v := range relations {
		result = append(result, v.UserId)
	}
	return result, nil
}

func getFollowerIdList(userId int64) ([]int64, error) {
	var relations []Relation
	err := dbConn.Where("user_id = ?", userId).Find(&relations).Error
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, v := range relations {
		result = append(result, v.FollowerId)
	}
	return result, nil
}

func addFollowRelationToRedis(followerId int64, followings []int64) {
	ctx := context.Background()

	for _, following := range followings {
		rdFollows.AddFollow(ctx, following, followerId)
	}
}

func addFollowerRelationToRedis(userId int64, followers []int64) {
	ctx := context.Background()

	for _, follower := range followers {
		rdFollows.AddFollower(ctx, userId, follower)
	}
}

// GetFriendIdList 获取用户好友列表
func GetFriendIdList(ctx context.Context, userId int64) ([]int64, error) {
	if !rdFollows.CheckFollow(ctx, userId) {
		following, err := getFollowerIdList(userId)
		if err != nil {
			return *new([]int64), err
		}
		addFollowerRelationToRedis(userId, following)
	}

	if !rdFollows.CheckFollower(ctx, userId) {
		followers, err := getFollowerIdList(userId)
		if err != nil {
			return *new([]int64), err
		}
		addFollowerRelationToRedis(userId, followers)
	}

	return rdFollows.GetFriend(ctx, userId), nil
}
