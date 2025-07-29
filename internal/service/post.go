package service

import (
	"qwqserver/internal/global"
	"qwqserver/internal/repository"
	postService "qwqserver/internal/service/post"
)

type PostService struct {
	postService.Service
}

// Create 创建文章
//func (s *PostService) Create() (res *global.HTTPResult) {
//	res = &global.HTTPResult{}
//
//	//
//	postRepo, err := repository.NewPostRepository()
//	if err != nil {
//		res.Code = 500
//		res.Message = "获取数据库连接失败: " + err.Error()
//		return
//	}
//
//	if err = postRepo.Create(context.Background(), s.Post); err != nil {
//		res.Code = 500
//		res.Message = "创建文章失败: " + err.Error()
//		return
//	}
//
//	res.Code = 200
//	res.Message = "创建文章成功"
//	res.Data = s
//	return
//}

// PostList 获取文章列表
func PostList(page, pageSize int) (res *global.HTTPResult) {
	res = &global.HTTPResult{}

	_, err := repository.NewPostRepository()
	if err != nil {
		res.Code = 500
		res.Message = "获取数据库连接失败: " + err.Error()
		return
	}
	panic("not implemented")
}
