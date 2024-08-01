package router

import (
	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/chat-kafka/internal/api/dto"
	"github.com/seedlings-calm/chat-kafka/internal/api/handler"
)

func SetUserRouter(v1 *gin.RouterGroup) {
	dto := new(dto.UserOper)
	r := v1.Group("/user", handler.Auth())
	{
		r.GET("/getGroups", dto.GetGroups)
	}
}
