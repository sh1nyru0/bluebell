package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	//注册路由
	v1.POST("/signup",controller.SignUpHandler)
	v1.POST("/login",controller.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件

	{
		v1.GET("/community",controller.CommunityHandler)
		v1.GET("/community/:id",controller.CommunityDetailHandler)
		v1.POST("/post",controller.CreatePostHandler)
		v1.GET("/post/:id",controller.GetPostHandler)
		v1.GET("/posts/",controller.GetPostListHandler)
		// 投票
		v1.POST("/vote",controller.PostVoteHandler)
	}


	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}

