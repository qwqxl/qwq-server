package handler

import (
	"qwqserver/internal/service"
	postService "qwqserver/internal/service/post"
	"qwqserver/pkg/httpcore"
)

type PostHandler struct{}

func (handle *PostHandler) Create(c httpcore.Context) {
	post := &postService.CreateRequest{}
	if err := c.ShouldBindJSON(post); err != nil {
		if post.Title != "" {
			httpcore.Fail(c, "title not found")
			return
		}
		if post.Content != "" {
			httpcore.Fail(c, "content not found")
			return
		}
		httpcore.Fail(c, "参数错误: "+err.Error())
		return
	}
	serv := service.PostService{}
	data, err := serv.Create(*post)
	if err != nil {
		httpcore.Fail(c, err.Error())
		return
	}
	httpcore.Success(c, data)
}
