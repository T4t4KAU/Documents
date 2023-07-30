package db

import (
	"douyin/biz/mw/redis"
	"douyin/pkg/constants"
	"gorm.io/gorm"
	"time"
)

var rdFollows redis.Follows

type Follows struct {
	ID         int64          `json:"id"`
	UserId     int64          `json:"user_id"`
	FollowerId int64          `json:"follower_id"`
	CreatedAt  time.Time      `json:"create_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

func (Follows) TableName() string {
	return constants.FollowsTableName
}

// AddNewFollow 添加关注信息
func AddNewFollow(follow *Follows) (bool, error) {
	err := dbConn.Create(follow).Error
	if err != nil {
		return false, err
	}

	if rdFollows.CheckFollow(follow.FollowerId) {
		rdFollows.AddFollow(follow.UserId, follow.FollowerId)
	}
	if rdFollows.CheckFollower(follow.UserId) {
		rdFollows.AddFollower(follow.UserId, follow.FollowerId)
	}

	return true, nil
}

// DeleteFollow 删除关注信息
func DeleteFollow(follow *Follows) (bool, error) {
	err := dbConn.Where("user_id = ? AND follower_id = ?",
		follow.UserId, follow.FollowerId).Error
	if err != nil {
		return false, err
	}

	if rdFollows.CheckFollower(follow.FollowerId) {
		rdFollows.DelFollower(follow.UserId, follow.FollowerId)
	}
	if rdFollows.CheckFollower(follow.UserId) {
		rdFollows.DelFollow(follow.UserId, follow.FollowerId)
	}

	return true, nil
}

// QueryFollowExist 查询关注是否存在
func QueryFollowExist(userId, followerId int64) (bool, error) {
	if rdFollows.CheckFollow(followerId) {
		return rdFollows.ExistFollow(userId, followerId), nil
	}
	if rdFollows.CheckFollower(userId) {
		return rdFollows.ExistFollower(userId, followerId), nil
	}

	follow := Follows{
		UserId:     userId,
		FollowerId: followerId,
	}

	err := dbConn.Where("user_id = ? AND follower_id = ?", userId, followerId).Find(&follow).Error
	if err != nil {
		return false, err
	}
	if follow.ID == 0 {
		return false, nil
	}
	return true, nil
}

// GetFollowCount 获取询用户关注数量
func GetFollowCount(followerId int64) (int64, error) {
	if rdFollows.CheckFollower(followerId) {
		return rdFollows.CountFollower(followerId)
	}

	followings, err := getFollowIdList(followerId)
	if err != nil {
		return 0, err
	}

	go addFollowRelationToRedis(followerId, followings)
	return int64(len(followings)), nil
}

// 更新关注
func addFollowRelationToRedis(followerId int64, followings []int64) {
	for _, following := range followings {
		rdFollows.AddFollower(following, followerId)
	}
}

// GetFollowerCount 获取一个用户的粉丝数量
func GetFollowerCount(userId int64) (int64, error) {
	if rdFollows.CheckFollow(userId) {
		return rdFollows.CountFollower(userId)
	}

	followers, err := getFollowIdList(userId)
	if err != nil {
		return 0, err
	}

	go addFollowRelationToRedis(userId, followers)
	return int64(len(followers)), nil
}

// 更新redis缓存
func addFollowerRelationToRedis(userId int64, followers []int64) {
	for _, follower := range followers {
		rdFollows.AddFollower(userId, follower)
	}
}

// 从数据库获取用户粉丝列表
func getFollowIdList(followerId int64) ([]int64, error) {
	var followActions []Follows
	err := dbConn.Where("follower_id = ?", followerId).Find(&followActions).Error
	if err != nil {
		return nil, err
	}

	var result []int64
	for _, v := range followActions {
		result = append(result, v.UserId)
	}
	return result, nil
}

// GetFollowIdList 获取用户粉丝列表
func GetFollowIdList(followerId int64) ([]int64, error) {
	if rdFollows.CheckFollow(followerId) {
		return rdFollows.GetFollow(followerId), nil
	}
	return getFollowIdList(followerId)
}

// GetFriendIdList 获取用户好友列表
func GetFriendIdList(userId int64) ([]int64, error) {
	if !rdFollows.CheckFollow(userId) {
		following, err := getFollowIdList(userId)
		if err != nil {
			return *new([]int64), err
		}
		addFollowRelationToRedis(userId, following)
	}

	if !rdFollows.CheckFollow(userId) {
		followers, err := getFollowIdList(userId)
		if err != nil {
			return *new([]int64), err
		}
		addFollowRelationToRedis(userId, followers)
	}

	return rdFollows.GetFriend(userId), nil
}
