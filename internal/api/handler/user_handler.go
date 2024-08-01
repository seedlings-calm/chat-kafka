package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("验证")
		ctx.Next()
	}
}
