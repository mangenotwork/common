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

func (ctx *GinCtx) APIOutPut(code int64, msg string, data interface{}) {
	ctx.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      code,
		Msg:       msg,
		Date:      data,
		TimeStamp: time.Now().Unix(),
	})
	return
}

func (ctx *GinCtx) APIOutPutError(c *gin.Context, code int64, msg string) {
	c.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      code,
		Msg:       msg,
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

// GetPostArgs 获取参数
func (ctx *GinCtx) GetPostArgs(c *gin.Context, obj interface{}) error {
	return c.BindJSON(obj)
}

// ResponseJson 统一接口输出
type ResponseJson struct {
	Code      int64       `json:"code"`
	Msg       string      `json:"msg"`
	Date      interface{} `json:"data"`
	TimeStamp int64       `json:"timeStamp"`
}

// APIOutPut 统一接口输出方法
func APIOutPut(c *gin.Context, code int64, msg string, data interface{}) {
	c.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      code,
		Msg:       msg,
		Date:      data,
		TimeStamp: time.Now().Unix(),
	})
	return
}

// APIOutPutError 统一接口输出错误方法
func APIOutPutError(c *gin.Context, code int64, msg string) {
	c.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      code,
		Msg:       msg,
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

// GetPostArgs 获取参数
func GetPostArgs(c *gin.Context, obj interface{}) error {
	return c.BindJSON(obj)
}
