package httpMiddleware

import (
	"net"
	"net/http"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
)

// 常用封装中间件

// GetIP 获取ip
// - X-Real-IP：只包含客户端机器的一个IP，如果为空，某些代理服务器（如Nginx）会填充此header。
// - X-Forwarded-For：一系列的IP地址列表，以,分隔，每个经过的代理服务器都会添加一个IP。
// - RemoteAddr：包含客户端的真实IP地址。 这是Web服务器从其接收连接并将响应发送到的实际物理IP地址。 但是，如果客户端通过代理连接，它将提供代理的IP地址。
//
// RemoteAddr是最可靠的，但是如果客户端位于代理之后或使用负载平衡器或反向代理服务器时，它将永远不会提供正确的IP地址，因此顺序是先是X-REAL-IP，
// 然后是X-FORWARDED-FOR，然后是 RemoteAddr。 请注意，恶意用户可以创建伪造的X-REAL-IP和X-FORWARDED-FOR标头。
func GetIP(r *http.Request) (ip string) {
	for _, ip := range strings.Split(r.Header.Get("X-Forward-For"), ",") {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	if ip = r.Header.Get("X-Real-IP"); net.ParseIP(ip) != nil {
		return ip
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return "0.0.0.0"
}

// CorsHandler 跨域中间件
func CorsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
			"Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,"+
			"Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		ctx.Header("Access-Control-Max-Age", "172800")          // 缓存请求信息 单位为秒
		ctx.Header("Access-Control-Allow-Credentials", "false") //  跨域请求是否需要带cookie信息 默认设置为true
		//ctx.Header("content-type", "application/json")          // 设置返回格式是json
		//Release all option pre-requests
		if ctx.Request.Method == http.MethodOptions {
			ctx.JSON(http.StatusOK, "Options Request!")
		}
		ctx.Next()
	}
}

// SetSockOptInt 启动多个web服务进程，做本地负载
func SetSockOptInt() net.ListenConfig {
	return net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var opErr error
			if err := c.Control(func(fd uintptr) {
				opErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			}); err != nil {
				return err
			}
			return opErr
		},
		KeepAlive: 0,
	}
}

// TODO 限流

// TODO Auth

// TODO 请求上报，请求记录

// TODO JWT验证

// TODO id限流

// TODO 权限验证
