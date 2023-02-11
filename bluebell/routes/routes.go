package routes

import (
	"net/http"
	"bluebell/logger"
	"bluebell/settings"
	"bluebell/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//注册路由
	r.POST("/signup",controller.SignUpHandler)
	r.POST("/login",controller.LoginHandler)

	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}