package handler

import (
	"github.com/gin-gonic/gin"
	"qwqserver/internal/model"
	"qwqserver/internal/service"
)

func PostCreate(c *gin.Context) (res *model.Result) {
	res = &model.Result{}
	post := &model.Post{}
	post.AuthorID = 1
	if err := c.ShouldBindJSON(post); err != nil || post.Title == "" || post.Content == "" {
		res.Code = 400
		res.Message = "参数错误: " + err.Error()
		return
	}
	return service.PostCreate(post)
}
