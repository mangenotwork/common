package ginHelper

import (
	"fmt"
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

func NewGinCtx() *GinCtx {
	return &GinCtx{
		&gin.Context{},
	}
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
		Code:      TokenInvalid,
		Msg:       "鉴权失败,请前往登录!",
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

// GetPostArgs 获取参数
func (ctx *GinCtx) GetPostArgs(obj interface{}) error {
	return ctx.Context.ShouldBindJSON(obj)
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
	if has >= number {
		has = number
	}

	for i := 0; i < has; i++ {
		if pg-i > pg-2 && pg-i > 1 {
			temp1 = temp1 - 1
			pgs = append(pgs, temp1)
		} else {
			pgs = append(pgs, temp2)
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
		Code:      TokenInvalid,
		Msg:       "鉴权失败,请前往登录!",
		Date:      "",
		TimeStamp: time.Now().Unix(),
	})
	return
}

func (ctx *GinCtx) GetIP() string {
	if ip, ok := ctx.Get(ReqIP); ok {
		return ip.(string)
	}
	return ""
}

func (ctx *GinCtx) GetLang() string {
	if lang, ok := ctx.Get(Lang); ok {
		return lang.(string)
	}
	return ""
}

func (ctx *GinCtx) GetSource() string {
	if source, ok := ctx.Get(Source); ok {
		return source.(string)
	}
	return ""
}
