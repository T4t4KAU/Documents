package service

import (
	"context"
	"douyin/biz/dal/db"
	"douyin/biz/model/basic/user"
	"douyin/biz/model/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

// UserRegister 用户注册服务
func (s *UserService) UserRegister(req *user.DouyinUserRegisterRequest) (int64, error) {
	u, err := db.QueryUserByName(req.Username)
	if err != nil {
		return 0, err
	}
	if *u != (db.User{}) {
		return 0, errno.UserAlreadyExistErr
	}

	hashedPassword, _ := utils.CryptPassword(req.Password)

	uid, err := db.CreateUser(&db.User{
		UserName:        req.Username,
		Password:        hashedPassword,
		Avatar:          constants.TestAva,        // 测试头像
		BackgroundImage: constants.TestBackground, // 测试背景
	})

	return uid, err
}

// UserInfo 获取用户信息
func (s *UserService) UserInfo(req *user.DouyinUserRequest) (*common.User, error) {
	currentId, exists := s.c.Get("current_user_id")
	if !exists {
		currentId = int64(0)
	}
	return s.GetUserInfo(req.UserId, currentId.(int64))
}

// GetUserInfo 传入查询用户ID和当前用户ID 获取用户信息
func (s *UserService) GetUserInfo(queryId, userId int64) (*common.User, error) {
	u := &common.User{}

	errChan := make(chan error, 7)
	defer close(errChan)

	var wg sync.WaitGroup
	wg.Add(7)

	go func() {
		dbUser, err := db.QueryUserById(queryId)
		if err != nil {
			errChan <- err
		} else {
			u.Name = dbUser.UserName
			u.Avatar = constants.TestAva
			u.BackgroundImage = constants.TestBackground
			u.Signature = dbUser.Signature
		}
		wg.Done()
	}()

	go func() {
		count, err := db.GetPublishCountById(queryId)
		if err != nil {
			errChan <- err
		} else {
			u.WorkCount = count
		}
		wg.Done()
	}()

	go func() {
		count, err := db.GetFollowerCount(queryId)
		if err != nil {
			errChan <- err
		} else {
			u.FollowerCount = count
		}
		wg.Done()
	}()

	go func() {
		count, err := db.GetFollowCount(queryId)
		if err != nil {
			errChan <- err
		} else {
			u.FollowCount = count
		}
		wg.Done()
	}()

	go func() {
		if userId != 0 {
			isFollow, err := db.QueryFollowExist(userId, queryId)
			if err != nil {
				errChan <- err
			} else {
				u.IsFollow = isFollow
			}
		} else {
			u.IsFollow = false
		}
		wg.Done()
	}()

	go func() {
		favoriteCount, err := db.GetFavoriteCountByUserId(queryId)
		if err != nil {
			errChan <- err
		} else {
			u.FavoriteCount = favoriteCount
		}
		wg.Done()
	}()

	go func() {
		totalFavorited, err := db.QueryTotalFavoritedByAuthorId(queryId)
		if err != nil {
			errChan <- err
		} else {
			u.TotalFavorited = totalFavorited
		}
		wg.Done()
	}()

	wg.Wait()
	select {
	case result := <-errChan:
		return &common.User{}, result
	default:
	}

	u.Id = queryId
	return u, nil
}
