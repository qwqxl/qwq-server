package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, PingResponse{
		Code: http.StatusOK,
		Msg:  "pong",
		Data: nil,
	})

}

type PingResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}
