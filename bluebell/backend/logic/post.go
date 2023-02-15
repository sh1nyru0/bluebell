package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error){
	// 1.生成post id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil{
		return err
	}
	err = redis.CreatePost(p.ID)
	return 
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

// GetPostList 获取帖子列表
func GetPostList(page,size int64) (data []*models.ApiPostDetail,err error){
	posts,err := mysql.GetPostList(page,size)
	if err != nil{
		return nil,err
	}
	data = make([]*models.ApiPostDetail,0,len(posts))
	for _,post := range posts {
		// 根据作者id查询作者信息
		user,err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id",post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详情信息
		community,err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil{
			zap.L().Error("mysql.GetUserById(post.AuthorID failed)",
				zap.Int64("community_id",post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			Post: post,
			CommunityDetail: community,
		}
		data = append(data,postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList)(data []*models.ApiPostDetail,err error){
	// 1.去redis查询id列表
	ids,err := redis.GetPostIDsInOrder(p)
	if err != nil{
		return
	}
	if len(ids) == 0{
		zap.L().Warn("GetPostIDsInOrder(p) return 0 data")
		return 
	}
	zap.L().Debug("GetPostList2",zap.Any("ids",ids))
	// 2.根据id去Mysql数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回
	posts,err := mysql.GetPostListByIDs(ids)
	if err != nil{
		return
	}
	// 将帖子的作者及分区信息查询出来填充到帖子中                                                                       
	for _,post := range posts {
		// 根据作者id查询作者信息
		user,err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id",post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详情信息
		community,err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil{
			zap.L().Error("mysql.GetUserById(post.AuthorID failed)",
				zap.Int64("community_id",post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			Post: post,
			CommunityDetail: community,
		}
		data = append(data,postDetail)
	}
	return
}