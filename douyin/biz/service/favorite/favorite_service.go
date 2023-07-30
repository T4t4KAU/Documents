package service

import (
	"context"
	"douyin/biz/dal/db"
	"douyin/biz/model/common"
	"douyin/biz/model/interact/favorite"
	service "douyin/biz/service/feed"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
)

type FavoriteService struct {
	ctx context.Context
	c   *app.RequestContext
}

// NewFavoriteService 创建点赞服务
func NewFavoriteService(ctx context.Context, c *app.RequestContext) *FavoriteService {
	return &FavoriteService{ctx: ctx, c: c}
}

// FavoriteAction 视频点赞
func (s *FavoriteService) FavoriteAction(req *favorite.DouyinFavoriteActionRequest) (bool, error) {
	_, err := db.CheckVideoExistById(req.VideoId)
	if err != nil {
		return false, err
	}
	if req.ActionType != constants.FavoriteActionType && req.ActionType != constants.UnFavoriteActionType {
		return false, errno.ParamErr
	}

	currentId, _ := s.c.Get("current_user_id")
	newFavoriteRelation := &db.Favorites{
		UserId:  currentId.(int64),
		VideoId: req.VideoId,
	}

	favoriteExist, _ := db.QueryFavoriteExist(newFavoriteRelation.UserId, newFavoriteRelation.VideoId)
	if req.ActionType == constants.FavoriteActionType {
		if favoriteExist {
			return false, errno.FavoriteRelationAlreadyExistErr
		}
		return db.AddNewFavorite(newFavoriteRelation)
	} else {
		if !favoriteExist {
			return false, errno.FavoriteRelationNotExistErr
		}
		return db.DeleteFavorite(newFavoriteRelation)
	}
}

// GetFavoriteList 获取点赞列表
func (s *FavoriteService) GetFavoriteList(req *favorite.DouyinFavoriteListRequest) ([]*common.Video, error) {
	queryId := req.UserId
	_, err := db.CheckUserExistById(queryId)
	if err != nil {
		return nil, err
	}

	currentId, _ := s.c.Get("current_user_id")
	videoIds, err := db.GetFavoriteIdList(queryId)
	dbVideos, err := db.GetVideosByVideoIdList(videoIds)
	if err != nil {
		return nil, err
	}

	var videos []*common.Video
	fs := service.NewFeedService(s.ctx, s.c)
	err = fs.CopyVideos(&videos, &dbVideos, currentId.(int64))
	favoriteList := make([]*common.Video, 0)

	for _, item := range videos {
		video := &common.Video{
			Id: item.Id,
			Author: &common.User{
				Id:              item.Author.Id,
				Name:            item.Author.Name,
				FollowCount:     item.Author.FollowCount,
				FollowerCount:   item.Author.FollowerCount,
				Avatar:          item.Author.Avatar,
				BackgroundImage: item.Author.BackgroundImage,
				Signature:       item.Author.Signature,
				TotalFavorited:  item.Author.TotalFavorited,
				WorkCount:       item.Author.WorkCount,
			},
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: item.FavoriteCount,
			CommentCount:  item.CommentCount,
			IsFavorite:    item.IsFavorite,
			Title:         item.Title,
		}
		favoriteList = append(favoriteList, video)
	}

	return favoriteList, err
}
