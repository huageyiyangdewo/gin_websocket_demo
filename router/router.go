package router

import (
	"chat/api"
	"chat/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSONP(200, "success")
		})
		v1.POST("register", api.UserRegister)
		v1.GET("ws", service.Handler)
	}

	return r
}
