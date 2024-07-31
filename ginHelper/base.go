package ginHelper

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func NewGinServer() {

}

func OutHtml() {

}

const (
	ReqIP  = "reqIp"
	Lang   = "lang"   // cn: 中文
	Source = "source" // 来源，web:1, h5:2, android:3, ios:4
)

const (
	SourceWeb     = "1"
	SourceH5      = "2"
	SourceAndroid = "3"
	SourceIos     = "4"
)

// ResponseJson 统一接口输出
type ResponseJson struct {
	Code      ResponseCode `json:"code"`
	Msg       string       `json:"msg"`
	Date      any          `json:"data"`
	TimeStamp int64        `json:"timestamp"`
}

// ResponseCode 统一接口输出码
type ResponseCode int64

const (
	SuccessCode      ResponseCode = 0    // 接口成功
	ErrorCode        ResponseCode = 1    // 接口内部错误(业务错误，msg显示具体信息)
	ParamViolation   ResponseCode = 2    // 参数不合法
	InterfaceInvalid ResponseCode = 3    // 接口无效
	ServerStop       ResponseCode = 4    // 服务已暂停
	DataNone         ResponseCode = 5    // 接口没查询到数据
	TokenInvalid     ResponseCode = 1001 // token无效
	TokenExpire      ResponseCode = 1002 // token已过期
	NoPermissions    ResponseCode = 1003 // 没有权限访问
)

// OutPut 统一接口输出方法
func OutPut(c *gin.Context, msg string, data interface{}) {
	c.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      SuccessCode,
		Msg:       msg,
		Date:      data,
		TimeStamp: time.Now().Unix(),
	})
	return
}

// OutPutError 统一接口输出错误方法
func OutPutError(c *gin.Context, msg string) {
	c.IndentedJSON(http.StatusOK, ResponseJson{
		Code:      ErrorCode,
		Msg:       msg,
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

// TokenInvalidOut 鉴权失败
func TokenInvalidOut(c *gin.Context) {
	c.IndentedJSON(http.StatusForbidden, ResponseJson{
		Code:      TokenInvalid,
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
