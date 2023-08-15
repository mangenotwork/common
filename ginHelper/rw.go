package ginHelper

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *GinCtx)

/*
r.GET("/test", ginHelper.Handle(controllers.Test))
*/

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &GinCtx{
			c,
		}
		h(ctx)
	}
}

// GinCtx 给gin context 扩展方法
type GinCtx struct {
	*gin.Context
}

func (ctx *GinCtx) APIOutPut(data interface{}, msg string) {
	if data == nil {
		data = ""
	}
	ctx.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      SuccessCode,
		Msg:       msg,
		Date:      data,
		TimeStamp: time.Now().Unix(),
	})
	return
}

func (ctx *GinCtx) APIOutPutError(err error, msg string) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	ctx.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      ErrorCode,
		Msg:       msg,
		Date:      errStr,
		TimeStamp: time.Now().Unix(),
	})
	return
}

// AuthErrorOut 鉴权失败
func (ctx *GinCtx) AuthErrorOut() {
	ctx.IndentedJSON(http.StatusForbidden, ResponseJson{
		Code:      AuthErrorCode,
		Msg:       "鉴权失败,请前往登录!",
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

// GetPostArgs 获取参数
func (ctx *GinCtx) GetPostArgs(obj interface{}) error {
	return ctx.Context.BindJSON(obj)
}

// ResponseJson 统一接口输出
type ResponseJson struct {
	Code      int64       `json:"code"` // succeed:0  err:1
	Msg       string      `json:"msg"`
	Date      interface{} `json:"data"`
	TimeStamp int64       `json:"timeStamp"`
}

const (
	SuccessCode   int64 = 0
	ErrorCode     int64 = 1
	AuthErrorCode int64 = 403
)

// APIOutPut 统一接口输出方法
func APIOutPut(c *gin.Context, msg string, data interface{}) {
	c.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      SuccessCode,
		Msg:       msg,
		Date:      data,
		TimeStamp: time.Now().Unix(),
	})
	return
}

// APIOutPutError 统一接口输出错误方法
func APIOutPutError(c *gin.Context, msg string) {
	c.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      ErrorCode,
		Msg:       msg,
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

// AuthErrorOut 鉴权失败
func AuthErrorOut(c *gin.Context) {
	c.IndentedJSON(http.StatusForbidden, ResponseJson{
		Code:      AuthErrorCode,
		Msg:       "鉴权失败,请前往登录!",
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

// GetPostArgs 获取参数
func GetPostArgs(c *gin.Context, obj interface{}) error {
	return c.BindJSON(obj)
}
