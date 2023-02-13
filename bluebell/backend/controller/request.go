package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrorUserNoLogin = errors.New("用户未登录")

// getCurrentUser获取当前登录的用户的ID
func getCurrentUserID(c *gin.Context) (userID int64,err error) {
	uid,ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNoLogin
		return
	}
	userID,ok = uid.(int64)
	if !ok {
		err = ErrorUserNoLogin
		return
	}
	return
}