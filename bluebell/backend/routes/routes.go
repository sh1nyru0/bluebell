package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/settings"
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

	v1.GET("/ping",middlewares.JWTAuthMiddleware(),func(c *gin.Context) {
		//如果是登录的用户,判断请求头中是否有 有效的JWT ？	
		c.String(http.StatusOK,"ping")
	})

	v1.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}

