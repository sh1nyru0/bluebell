package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvaildParam
	CodeUserExist
	CodeUserNoExist
	CodeInvaildPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvaildToken	
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvaildParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNoExist:     "用户名不存在",
	CodeInvaildPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeNeedLogin: "需要登录",
	CodeInvaildToken: "无效的token",
	
}

func (c ResCode) Msg() string{
	msg,ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}