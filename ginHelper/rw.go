package ginHelper

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"net/http"
	"sort"
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

func (ctx *GinCtx) GetParam(key string) string {
	return ctx.Context.Param(key)
}

func (ctx *GinCtx) GetParamInt(key string) int {
	v := ctx.Context.Param(key)
	return utils.AnyToInt(v)
}

func (ctx *GinCtx) GetParamInt64(key string) int64 {
	v := ctx.Context.Param(key)
	return utils.AnyToInt64(v)
}

func (ctx *GinCtx) GetQuery(key string) string {
	v, _ := ctx.Context.GetQuery(key)
	return v
}

func (ctx *GinCtx) GetQueryInt(key string) int {
	v, _ := ctx.Context.GetQuery(key)
	return utils.AnyToInt(v)
}

func (ctx *GinCtx) GetQueryInt64(key string) int64 {
	v, _ := ctx.Context.GetQuery(key)
	return utils.AnyToInt64(v)
}

type Page struct {
	Number int
	Action bool
	Url    string
}

func (ctx *GinCtx) PageListInt(pg, number, count, size int) []int {
	pgs := make([]int, 0)
	temp1 := pg
	temp2 := pg

	has := (count / size) + 1
	log.Info("has = ", has)
	if has >= number {
		has = number
	}

	for i := 0; i < has; i++ {
		if int(pg)-i > int(pg)-2 && int(pg)-i > 1 {
			temp1 = temp1 - 1
			pgs = append(pgs, int(temp1))
		} else {
			pgs = append(pgs, int(temp2))
			temp2 = temp2 + 1
		}
	}
	sort.Ints(pgs)
	return pgs
}

// PageList urlTemp = /a?pg=%d
func (ctx *GinCtx) PageList(pg, number, count, size int, urlTemp string) []*Page {
	pgIntList := ctx.PageListInt(pg, number, count, size)
	list := make([]*Page, 0)
	for _, i := range pgIntList {
		p := &Page{
			Number: i,
			Action: false,
		}
		if p.Number == pg {
			p.Action = true
		}
		if len(urlTemp) > 0 {
			p.Url = fmt.Sprintf(urlTemp, p.Number)
		}
		list = append(list, p)
	}
	return list
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
