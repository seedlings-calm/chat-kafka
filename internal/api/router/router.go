package router

import "github.com/gin-gonic/gin"

func InitEngine() (r *gin.Engine) {
	r = gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	g := r.Group("v1")
	SetUserRouter(g)
	return r
}
