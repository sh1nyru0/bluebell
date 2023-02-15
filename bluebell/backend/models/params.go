package models

//定义请求的参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票数据
type ParamVoteData struct {
	// UserID 从当前的请求中获取当前的用户
	PostID    string `json:"post_id" binging:"required"`    // 帖子id
	Direction int8   `json:"direcation,string" binging:"oneof=1 0 -1"` // 赞成票(1)还是反对票(-1)取消投票(0)
}
