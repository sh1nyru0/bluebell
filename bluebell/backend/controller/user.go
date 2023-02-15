package controller

import (
	"errors"
	"fmt"

	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignupHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c,CodeInvaildParam)
			return
		}
		ResponseErrorWithMsg(c,CodeInvaildParam,removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务处理
	if err := logic.SignUp(p);err != nil{
		zap.L().Error("logic.SignUp failer",zap.Error(err))
		if errors.Is(err,mysql.ErrorUserExist){
			ResponseError(c,CodeUserExist)
			return
		}
		ResponseError(c,CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c,nil)
}

func LoginHandler(c *gin.Context){
	// 1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p);err != nil{
		//请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断err是不是validationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c,CodeInvaildParam)
			return
		}
		ResponseErrorWithMsg(c,CodeInvaildParam,removeTopStruct(errs.Translate(trans)))
		return
	}
	
	// 2.业务逻辑处理
	user,err := logic.Login(p)
	if err != nil{
		zap.L().Error("Login with invalid param",zap.String("username",p.Username),zap.Error(err))
		if errors.Is(err,mysql.ErrorUserNoExist){
			ResponseError(c,CodeInvalidPassword)
			return
		}
		ResponseError(c,CodeInvalidPassword)
		return
	}
	// 3.返回响应
	ResponseSuccess(c,gin.H{
		"user_id": fmt.Sprintf("%d",user.UserID),//id值大于(1<<53)-1 int64
		"user_name": user.Username,
		"token": user.Token,
	})
}