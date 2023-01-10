package ginHelper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

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
