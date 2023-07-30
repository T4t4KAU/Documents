package service

import (
	"context"
	"douyin/biz/dal/db"
	"douyin/biz/model/common"
	"douyin/biz/model/interact/comment"
	service "douyin/biz/service/user"
	"douyin/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
)

type CommentService struct {
	ctx context.Context
	c   *app.RequestContext
}

// AddNewComment 添加一条评论并返回
func (c *CommentService) AddNewComment(req *comment.DouyinCommentActionRequest) (*comment.Comment, error) {
	currentId, _ := c.c.Get("current_user_id")
	videoId := req.VideoId
	actionType := req.ActionType
	commentText := req.CommentText
	commentId := req.CommentId
	cmt := &comment.Comment{}

	// 发表评论
	if actionType == 1 {
		dbComment := &db.Comment{
			UserId:      currentId.(int64),
			VideoId:     videoId,
			CommentText: commentText,
		}
		err := db.AddNewComment(dbComment)
		if err != nil {
			return cmt, err
		}
		cmt.Id = dbComment.ID
		cmt.CreateDate = dbComment.CreatedAt.Format("01-02")
		cmt.Content = dbComment.CommentText
		cmt.User, err = c.getUserInfoById(currentId.(int64), currentId.(int64))
		if err != nil {
			return cmt, err
		}
	} else {
		// 删除评论
		err := db.DeleteCommentById(commentId)
		if err != nil {
			return cmt, err
		}
	}

	return cmt, nil
}

func (c *CommentService) getUserInfoById(currentId, uid int64) (*common.User, error) {
	u, err := service.NewUserService(c.ctx, c.c).GetUserInfo(uid, currentId)
	if err != nil {
		return nil, err
	}

	commentUser := &common.User{
		Id:              u.Id,
		Name:            u.Name,
		FollowCount:     u.FollowCount,
		FollowerCount:   u.FollowerCount,
		IsFollow:        u.IsFollow,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
		TotalFavorited:  u.TotalFavorited,
		WorkCount:       u.WorkCount,
		FavoriteCount:   u.FavoriteCount,
	}

	return commentUser, nil
}

// CommentList 返回指定视频的评论列表
func (c *CommentService) CommentList(req *comment.DouyinCommentListRequest) (*comment.DouyinCommentListResponse, error) {
	resp := &comment.DouyinCommentListResponse{}
	videoId := req.VideoId

	currentId, _ := c.c.Get("current_user_id")
	dbComments, err := db.GetCommentListByVideoId(videoId)
	if err != nil {
		return resp, err
	}

	var comments []*comment.Comment
	err = c.copyComment(&comments, &dbComments, currentId.(int64))
	if err != nil {
		return resp, err
	}

	resp.CommentList = comments
	resp.StatusMsg = errno.SuccessMsg
	resp.StatusCode = errno.SuccessCode

	return resp, nil
}

func (c *CommentService) copyComment(result *[]*comment.Comment, data *[]*db.Comment, currentId int64) error {
	for _, item := range *data {
		cmt := c.createComment(item, currentId)
		*result = append(*result, cmt)
	}
	return nil
}

func (c *CommentService) createComment(data *db.Comment, uid int64) *comment.Comment {
	cmt := &comment.Comment{
		Id:         data.ID,
		Content:    data.CommentText,
		CreateDate: data.CreatedAt.Format("01-02"),
	}

	userInfo, err := c.getUserInfoById(uid, data.UserId)
	if err != nil {
		log.Printf("func error")
	}
	cmt.User = userInfo
	return cmt
}
