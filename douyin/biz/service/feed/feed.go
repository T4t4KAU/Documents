package service

import (
	"context"
	"douyin/biz/dal/db"
	"douyin/biz/model/basic/feed"
	"douyin/biz/model/common"
	service "douyin/biz/service/user"
	"douyin/pkg/constants"
	"douyin/pkg/utils"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"sync"
	"time"
)

type FeedService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewFeedService(ctx context.Context, c *app.RequestContext) *FeedService {
	return &FeedService{ctx: ctx, c: c}
}

// Feed 获取截止日期前的10个视频
func (s *FeedService) Feed(req *feed.DouyinFeedRequest) (*feed.DouyinFeedResponse, error) {
	resp := &feed.DouyinFeedResponse{}
	var lastTime time.Time

	if req.LatestTime == 0 {
		lastTime = time.Now()
	} else {
		lastTime = time.Unix(req.LatestTime/1000, 0)
	}

	fmt.Printf("LastTime: %v\n", lastTime)
	currentId, exists := s.c.Get("current_user_id")
	if !exists {
		currentId = int64(0)
	}

	dbVideos, err := db.GetVideosByLastTime(lastTime)
	if err != nil {
		return resp, err
	}

	videos := make([]*common.Video, 0, constants.VideoFeedCount)
	err = s.CopyVideos(&videos, &dbVideos, currentId.(int64))
	if err != nil {
		return resp, nil
	}
	resp.VideoList = videos
	if len(dbVideos) != 0 {
		resp.NextTime = dbVideos[len(dbVideos)-1].PublishTime.Unix()
	}
	return resp, nil
}

func (s *FeedService) CopyVideos(result *[]*common.Video, data *[]*db.Video, userId int64) error {
	for _, item := range *data {
		video := s.createVideo(item, userId)
		*result = append(*result, video)
	}
	return nil
}

func (s *FeedService) createVideo(data *db.Video, userId int64) *common.Video {
	video := &common.Video{
		Id:       data.ID,
		PlayUrl:  utils.URLConvert(s.ctx, s.c, data.PlayURL),
		CoverUrl: utils.URLConvert(s.ctx, s.c, data.CoverURL),
		Title:    data.Title,
	}

	var wg sync.WaitGroup
	wg.Add(4)

	// Get author information
	go func() {
		author, err := service.NewUserService(s.ctx, s.c).GetUserInfo(data.AuthorID, userId)
		if err != nil {
			log.Printf("GetUserInfo func error:" + err.Error())
		}
		video.Author = &common.User{
			Id:              author.Id,
			Name:            author.Name,
			FollowCount:     author.FollowCount,
			FollowerCount:   author.FollowerCount,
			IsFollow:        author.IsFollow,
			Avatar:          author.Avatar,
			BackgroundImage: author.BackgroundImage,
			Signature:       author.Signature,
			TotalFavorited:  author.TotalFavorited,
			WorkCount:       author.WorkCount,
			FavoriteCount:   author.FavoriteCount,
		}

		wg.Done()
	}()

	// 获取点赞数
	go func() {
		err := *new(error)
		video.FavoriteCount, err = db.GetFavoriteCount(data.ID)
		if err != nil {
			log.Printf("GetFavoriteCount func error:" + err.Error())
		}
		wg.Done()
	}()

	// 获取评论数
	go func() {
		err := *new(error)
		video.CommentCount, err = db.GetCommentCountByVideoId(data.ID)
		if err != nil {
			log.Printf("GetCommentCountByVideoID func error:" + err.Error())
		}
		wg.Done()
	}()

	// 是否点赞
	go func() {
		err := *new(error)
		video.IsFavorite, err = db.QueryFavoriteExist(userId, data.ID)
		if err != nil {
			log.Printf("QueryFavoriteExist func error:" + err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
	return video
}
