package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler创建帖子的处理函数
func CreatePostHandler(c *gin.Context){
	//1.获取参数及参数的校验

	p := new(models.Post)
	if err := c.ShouldBindJSON(p);err != nil{
		zap.L().Debug("c.ShouldBindJSON(p) error",zap.Any("",err))
		zap.L().Error("create post with invalid param")
		ResponseError(c,CodeInvaildParam)
		return
	}

	// 从c取到当前发请求的用户的ID
	userID,err := getCurrentUserID(c)
	if err != nil{
		ResponseError(c,CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2.创建帖子
 	if err := logic.CreatePost(p);err != nil{
		zap.L().Error("logic.CreatePost(p) failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c,nil)
}

// GetPostHandler获取帖子详情的处理函数
func GetPostHandler(c *gin.Context){
	// 1.从URL中获取参数(帖子的id)
	pidStr := c.Param("id")
	pid,err := strconv.ParseInt(pidStr,10,64)
	if err != nil{
		zap.L().Error("get post detail with invalid param",zap.Error(err))
		ResponseError(c,CodeInvaildParam)
		return
	}
	// 2.根据id取出帖子数据

	data,err := logic.GetPostById(pid)
	if err != nil{
		zap.L().Error("logic.GetPostById(pid) failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c,data)
}