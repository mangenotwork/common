package ginHelper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/mangenotwork/common/log"
	"html/template"
	"net/http"
)

// 使用 https://github.com/gorilla/csrf

var (
	CsrfName    = "cToken"
	CsrfAuthKey = "123456789"
)

// CSRFMiddleware
// use:
// ginHelper.CsrfName = "CsrfName"
// ginHelper.CsrfAuthKey = "CsrfAuthKey"
// Router.Use(ginHelper.CSRFMiddleware())
func CSRFMiddleware() gin.HandlerFunc {
	csrfMiddleware := csrf.Protect(
		[]byte(CsrfAuthKey),
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.CookieName(CsrfName),
		csrf.FieldName(CsrfName),
		csrf.RequestHeader(CsrfName),
		csrf.ErrorHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			log.Info(request.URL.Path)
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte(`非法请求!`))
		})),
	)
	// 这里使用adpater包将csrfMiddleware转换成gin的中间件返回值
	return adapter.Wrap(csrfMiddleware)
}

func CSRFTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		log.Info("token = ", token)
		c.Header(CsrfName, token)
	}
}

func SetCSRFTokenMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		c.Header(key, token)
		c.SetCookie(key, token, 3600*24, "/", "", false, false)
	}
}

// FormSetCSRF
// use:
//
//	c.HTML(http.StatusOK, "login.html", gin.H{
//		"csrf": ginHelper.FormSetCSRF(c.Request),
//	})
func FormSetCSRF(r *http.Request) template.HTML {
	fragment := fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`,
		CsrfName, csrf.Token(r))
	return template.HTML(fragment)
}
