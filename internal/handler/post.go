package handler

import (
	"github.com/gin-gonic/gin"
	"qwqserver/internal/common"
	"qwqserver/internal/model"
	"qwqserver/internal/service"
)

type PostHandler struct {
	*HandleBaseImpl
}

func NewPost() *PostHandler {
	return &PostHandler{}
}

func (handle *PostHandler) Create(c *gin.Context) (res *common.HTTPResult) {
	res = &common.HTTPResult{}
	post := &model.Post{}
	post.AuthorID = 1
	if err := c.ShouldBindJSON(post); err != nil || post.Title == "" || post.Content == "" {
		res.Code = 400
		res.Msg = "参数错误: "
		if err != nil {
			res.Msg += err.Error()
		}
		return
	}
	serv := service.NewPost(post)
	return serv.Create()
}
