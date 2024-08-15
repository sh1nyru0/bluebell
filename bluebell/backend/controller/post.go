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

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context){
	// 获取分页参数
	page,size := getPageInfo(c)
	// 获取数据
	data,err := logic.GetPostList(page,size)
	if err != nil{
		zap.L().Error("logic.GetPostList() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
	}
	ResponseSuccess(c,data)
	// 返回响应
}

// GetPostListHandler2升级版帖子列表接口
// 根据前端传来的参数动态获取帖子列表
// 按创建时间排序 或者 按照分数排序
// 1.获取参数
// 2.去redis查询id列表
// 3.根据id去数据库查询帖子详细信息
func GetPostListHandler2(c *gin.Context){
	// GET请求的query string参数(query string)：/api/v1/posts2?page=1&size=10&order=time
	// 初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page: 1,
		Size: 10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p);err != nil{
		zap.L().Error("GetPostListHandler2 with invalid params",zap.Error(err))
		ResponseError(c,CodeInvaildParam)
		return
	}
	data,err := logic.GetPostListNew(p) // 更新： 合二为一
	// 获取数据
	if err != nil{
		zap.L().Error("logic.GetPostList() failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c,data)
	// 返回响应
}