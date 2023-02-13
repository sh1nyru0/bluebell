package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error){
	// 1.生成post id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	return mysql.CreatePost(p)
	// 3.返回
}

// GetPostById 根据帖子id查询帖子详情数据
func GetPostById(pid int64) (data *models.ApiPostDetail,err error){
	// 查询数据并组合我们结构想用的数据
	post,err := mysql.GetPostById(pid)
	if err != nil{
		zap.L().Error("mysql.GetPostById(pid) failed",zap.Int64("pid",pid),zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user,err := mysql.GetUserById(post.AuthorID)
	if err != nil{
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
		zap.Int64("author_id",post.AuthorID),
		zap.Error(err))
		return
	}
	// 根据社区id查询社区详情信息
	community,err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil{
		zap.L().Error("mysql.GetUserById(post.AuthorID failed)",
			zap.Int64("community_id",post.CommunityID),
			zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:  user.Username,
		Post: post,
		CommunityDetail: community,
	}
	return 
}