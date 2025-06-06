package service

import (
	"context"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
)

// PostCreate 创建文章
func PostCreate(postRequest *model.Post) (res *model.Result) {
	res = &model.Result{}

	//
	postRepo, err := repository.NewPostRepository()
	if err != nil {
		res.Code = 500
		res.Message = "获取数据库连接失败: " + err.Error()
		return
	}

	if err = postRepo.Create(context.Background(), postRequest); err != nil {
		res.Code = 500
		res.Message = "创建文章失败: " + err.Error()
		return
	}

	res.Code = 200
	res.Message = "创建文章成功"
	res.Data = postRequest
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
