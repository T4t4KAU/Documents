package service

import (
	"context"
	"douyin/cmd/relation/dal/db"
	"douyin/kitex_gen/relation"
	"douyin/pkg/errno"
)

const (
	FOLLOW   = 1
	UNFOLLOW = 2
)

type RelationService struct {
	ctx context.Context
}

// NewRelationService 创建用户关系服务
func NewRelationService(ctx context.Context) *RelationService {
	return &RelationService{
		ctx: ctx,
	}
}

func (s *RelationService) RelationAction(req *relation.RelationActionRequest) (bool, error) {
	if req.ActionType != FOLLOW && req.ActionType != UNFOLLOW {
		return false, errno.ParamErr
	}

	if req.ToUserId == req.CurrentUserId {
		return false, errno.ParamErr
	}

	r := &db.Relation{
		UserId:     req.ToUserId,
		FollowerId: req.CurrentUserId,
	}

	exist, _ := db.CheckRelationExist(s.ctx, r.UserId, r.FollowerId)
	if req.ActionType == FOLLOW {
		if exist {
			return false, errno.FollowRelationAlreadyExistErr
		}
		return db.AddNewRelation(r)
	} else {
		if !exist {
			return false, errno.FollowRelationNotExistErr
		}
		return db.DeleteRelation(r)
	}
}
