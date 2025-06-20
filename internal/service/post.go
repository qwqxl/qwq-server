package service

import (
	"context"
	"qwqserver/internal/common"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
)

type PostService struct {
	*model.Post
}

func NewPost(post *model.Post) *PostService {
	return &PostService{
		Post: post,
	}
}

// Create 创建文章
func (s *PostService) Create() (res *common.HTTPResult) {
	res = &common.HTTPResult{}

	//
	postRepo, err := repository.NewPostRepository()
	if err != nil {
		res.Code = 500
		res.Msg = "获取数据库连接失败: " + err.Error()
		return
	}

	if err = postRepo.Create(context.Background(), s.Post); err != nil {
		res.Code = 500
		res.Msg = "创建文章失败: " + err.Error()
		return
	}

	res.Code = 200
	res.Msg = "创建文章成功"
	res.Data = s
	return
}

// PostList 获取文章列表
func PostList(page, pageSize int) (res *model.Result) {
	res = &model.Result{}

	_, err := repository.NewPostRepository()
	if err != nil {
		res.Code = 500
		res.Message = "获取数据库连接失败: " + err.Error()
		return
	}
	panic("not implemented")
}
