/*
 基于 golang net/http库 实现的http框架，无非golang官方三方引用
 实例见 _examples/mange_server
*/

package httpMiddleware

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/netutil"
	"golang.org/x/time/rate"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

const Logo = `

 █▄ ▄█ ▄▀▄ █▄ █   ▄▀▀ ██▀
 █ ▀ █ █▀█ █ ▀█   ▀▄█ █▄▄  v0.1.14
 https://github.com/mangenotwork/common

`

type Engine struct {
	mux  *http.ServeMux
	base func(next http.Handler) http.Handler
}

func NewEngine() *Engine {
	engine := &Engine{
		mux:  http.NewServeMux(),
		base: Base,
	}
	engine.mux.Handle("/hello", engine.base(http.HandlerFunc(Hello)))
	engine.mux.Handle("/health", engine.base(http.HandlerFunc(Health)))
	return engine
}

func SimpleEngine() *Engine {
	return &Engine{
		mux:  http.NewServeMux(),
		base: Base,
	}
}

func (engine *Engine) GetMux() *http.ServeMux {
	return engine.mux
}

func (engine *Engine) Router(path string, f func(w http.ResponseWriter, r *http.Request)) {
	engine.mux.Handle(path, engine.base(http.HandlerFunc(f)))
}

func (engine *Engine) RouterFunc(path string, f func(w http.ResponseWriter, r *http.Request)) {
	engine.mux.HandleFunc(path, f)
}

func (engine *Engine) Run(port string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	server := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    4 * time.Second,
		WriteTimeout:   4 * time.Second,
		IdleTimeout:    4 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        engine.mux,
		// tls.Config 有个属性 Certificates []Certificate
		// Certificate 里有属性 Certificate PrivateKey 分别保存 certFile keyFile 证书的内容
	}

	// 如果在高频高并发的场景下, 有很多请求是可以复用的时候
	// 最好开启 keep-alives 减少三次握手 tcp 销毁连接时有个 timewait 时间
	server.SetKeepAlivesEnabled(true)
	l, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Panic("Listen Err : %v", err)
		return
	}
	defer l.Close()

	// 开启最高连接数， 注意: linux/uinx有效果， win无效
	var rLimit syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("rLimit.Cur = ", rLimit.Cur)
	log.Info("rLimit.Max = ", rLimit.Max)
	rLimit.Cur = rLimit.Max
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Starting http server port -> ", port)
	// 对连接数的保护， 设置为最高连接数是 本机的最高连接数
	// https://github.com/golang/net/blob/master/netutil/listen.go
	l = netutil.LimitListener(l, int(rLimit.Max)*10)
	err = server.Serve(l)
	if err != nil {
		log.Panic("ListenAndServe err : ", err)
	}
}

func (engine *Engine) OpenPprof() {
	engine.mux.Handle("/debug/pprof/", engine.base(http.HandlerFunc(pprof.Index)))
	engine.mux.Handle("/debug/pprof/cmdline", engine.base(http.HandlerFunc(pprof.Cmdline)))
	engine.mux.Handle("/debug/pprof/profile", engine.base(http.HandlerFunc(pprof.Profile)))
	engine.mux.Handle("/debug/pprof/symbol", engine.base(http.HandlerFunc(pprof.Symbol)))
	engine.mux.Handle("/debug/pprof/trace", engine.base(http.HandlerFunc(pprof.Trace)))
}

var Path, _ = os.Getwd()

func Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("Hello l'm ManGe.\n" + Logo))
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("true"))
}

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (lrw *ResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Base Http基础中间件,日志
func Base(next http.Handler) http.Handler {
	return BaseFunc(next)
}

// BaseFunc Base Http基础中间件,日志
func BaseFunc(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
			// 中间件 上下文传递值
			data := map[string]interface{}{
			   "1": "one",
			   "2": "two",
			}
			ctx := context.WithValue(r.Context(), "data", data)
			r.WithContext(ctx)

			// 下文读值
			data := r.Context().Value("data").(ContextValue)["2"]
			fmt.Println(data) // 会打印 two
		*/
		start := time.Now().UnixNano()
		ip := GetIP(r)
		newW := NewResponseWriter(w)
		next.ServeHTTP(newW, r)
		logStr := fmt.Sprintf("%s \t| %s \t| %d \t| %f ms \t| %s \t", ip, r.Method, newW.StatusCode, float64(time.Now().UnixNano()-start)/100000, r.URL.String())
		log.HttpInfo(logStr)
	}
}

// ReqLimit 基础中间件 IP限流, IP黑白名单
func ReqLimit(ipv *IpVisitor, nextHeader http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip := GetIP(r)
		if ipv.IsBlackList(ip) {
			_, _ = w.Write([]byte("已经拉入黑名单，禁止访问！"))
			return
		}
		if !ipv.IsWhiteList(ip) {
			limiter := ipv.GetVisitor(ip)
			if limiter.AllowN(time.Now(), 1) == false {
				log.Info("ip限流")
				http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
				return
			}
		}
		nextHeader.ServeHTTP(w, r)
		log.Info("[%s] %s %s %v", ip, r.Method, r.URL.Path, time.Since(start))
	})
}

type IpVisitor struct {
	ips       map[string]*visitor
	mtx       sync.Mutex
	BlackList map[string]struct{}
	WhiteList map[string]struct{}
}

func NewIpVisitor() *IpVisitor {
	return &IpVisitor{
		ips: make(map[string]*visitor),
	}
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// CleanupVisitors 启动一个协成  10分钟查一下ip限流数据，看看有没有超过1小时删除记录，有就删除
// 主要目的的为了释放内存空间
func (ipv *IpVisitor) CleanupVisitors() {
	go func() {
		timer1 := time.NewTicker(10 * time.Millisecond)
		select {
		case <-timer1.C:
			for ip, v := range ipv.ips {
				if time.Now().Sub(v.lastSeen) > 1*time.Hour {
					ipv.mtx.Lock()
					delete(ipv.ips, ip)
					ipv.mtx.Unlock()
				}
			}
		}
	}()
}

func (ipv *IpVisitor) GetVisitor(ip string) *rate.Limiter {
	ipv.mtx.Lock()
	defer ipv.mtx.Unlock()
	v, exists := ipv.ips[ip]
	if !exists {
		return ipv.AddVisitor(ip)
	}
	// 更新时间
	v.lastSeen = time.Now()
	return v.limiter
}

func (ipv *IpVisitor) AddVisitor(ip string) *rate.Limiter {
	r := rate.Every(10 * time.Second) // 每*s往桶中放一个Token
	// 第一个参数是r Limit 代表每秒可以向Token桶中产生多少token
	// 第二个参数是b int b代表Token桶的容量大小
	limiter := rate.NewLimiter(r, 10)
	ipv.ips[ip] = &visitor{limiter, time.Now()}
	return limiter
}

func (ipv *IpVisitor) AddWhiteList(ip string) {
	ipv.WhiteList[ip] = struct{}{}
}

func (ipv *IpVisitor) IsWhiteList(ip string) (ok bool) {
	_, ok = ipv.WhiteList[ip]
	return
}

func (ipv *IpVisitor) DelWhiteList(ip string) {
	delete(ipv.WhiteList, ip)
}

func (ipv *IpVisitor) AddBlackList(ip string) {
	ipv.BlackList[ip] = struct{}{}
}

func (ipv *IpVisitor) IsBlackList(ip string) (ok bool) {
	_, ok = ipv.BlackList[ip]
	return
}

func (ipv *IpVisitor) DelBlackList(ip string) {
	delete(ipv.BlackList, ip)
}

// HttpOutBody HTTP  输出 json body 定义
// [Code]
// - 0 成功
// - 1001 参数错误
// - 2001 程序错误
type HttpOutBody struct {
	Code      int         `json:"code"`
	Timestamp int64       `json:"timestamp"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
}

const (
	BodyJSON       = "application/json; charset=utf-8"
	BodyAsciiJSON  = "application/json"
	BodyHTML       = "text/html; charset=utf-8"
	BodyJavaScript = "application/javascript; charset=utf-8"
	BodyXML        = "application/xml; charset=utf-8"
	BodyPlain      = "text/plain; charset=utf-8"
	BodyYAML       = "application/x-yaml; charset=utf-8"
	BodyDownload   = "application/octet-stream; charset=utf-8"
	BodyPDF        = "application/pdf"
	BodyJPG        = "image/jpeg"
	BodyPNG        = "image/png"
	BodyGif        = "image/gif"
	BodyWord       = "application/msword"
	BodyOctet      = "application/octet-stream"
)

func OutSucceedBodyJsonP(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", BodyJavaScript)
	body := &HttpOutBody{
		Code:      0,
		Timestamp: time.Now().Unix(),
		Msg:       "succeed",
		Data:      data,
	}
	bodyJson, err := body.JsonStr()
	if err != nil {
		OutErrBody(w, 2001, err)
	}
	_, _ = fmt.Fprintln(w, bodyJson)
	return
}

func OutSucceedBody(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", BodyJSON)
	body := &HttpOutBody{
		Code:      0,
		Timestamp: time.Now().Unix(),
		Msg:       "succeed",
		Data:      data,
	}
	bodyJson, err := body.JsonStr()
	if err != nil {
		OutErrBody(w, 2001, err)
	}
	_, _ = fmt.Fprintln(w, bodyJson)
	return
}

func OutErrBody(w http.ResponseWriter, code int, err error) {
	body := &HttpOutBody{
		Code:      code,
		Timestamp: time.Now().Unix(),
		Msg:       err.Error(),
		Data:      nil,
	}
	bodyJson, _ := body.JsonStr()
	_, _ = fmt.Fprintln(w, bodyJson)
	return
}

// OutStaticFile 输出静态文件
func OutStaticFile(w http.ResponseWriter, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprintln(w, err)
		return
	}
	w.Header().Add("Content-Type", BodyHTML)
	_, _ = fmt.Fprintln(w, string(data))
	return
}

func OutPdf(w http.ResponseWriter, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprintln(w, err)
		return
	}
	w.Header().Add("Content-Type", BodyPDF)
	_, _ = fmt.Fprintln(w, string(data))
	return
}

func OutJPG(w http.ResponseWriter, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprintln(w, err)
		return
	}
	w.Header().Add("Content-Type", BodyJPG)
	_, _ = fmt.Fprintln(w, string(data))
	return
}

// OutUploadFile 给客户端下载的静态文件
func OutUploadFile(w http.ResponseWriter, path, fileName string) {
	file, _ := os.Open(path)
	defer file.Close()

	fileHeader := make([]byte, 512)
	_, err := file.Read(fileHeader)
	if err != nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprintln(w, err)
		return
	}
	fileStat, _ := file.Stat()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	_, err = file.Seek(0, 0)
	if err != nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprintln(w, err)
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprintln(w, err)
		return
	}

}

func Out404(w http.ResponseWriter) {
	w.WriteHeader(404)
	_, _ = fmt.Fprintln(w, "404")
}

func (m *HttpOutBody) JsonStr() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Error("Umarshal failed:", err)
		return "", err
	}
	return string(b), nil
}

// GetUrlArg 获取URL的GET参数
func GetUrlArg(r *http.Request, name string) string {
	var arg string
	values := r.URL.Query()
	arg = values.Get(name)
	return arg
}

func GetUrlArgInt64(r *http.Request, name string) int64 {
	var arg string
	values := r.URL.Query()
	arg = values.Get(name)
	return utils.AnyToInt64(arg)
}

func GetUrlArgInt(r *http.Request, name string) int {
	var arg string
	values := r.URL.Query()
	arg = values.Get(name)
	return utils.AnyToInt(arg)
}

func GetJsonParam(r *http.Request, param interface{}) {
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&param)
}

func GetFromArg(r *http.Request, name string) string {
	return r.FormValue(name)
}

func GetFromFile(r *http.Request, name string) (multipart.File, *multipart.FileHeader, error) {
	return r.FormFile(name)
}

func GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(name)
}

func GetCookieVal(r *http.Request, name string) string {
	cookie, err := GetCookie(r, name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func SetCookie(w http.ResponseWriter, name, value string, t int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Expires:  time.Now().Add(time.Duration(t) * time.Second),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func SetCookieMap(w http.ResponseWriter, data map[string]string, t int) {
	for k, v := range data {
		SetCookie(w, k, v, t)
	}
}

func GetClientIp(r *http.Request) string {
	return GetIP(r)
}

func GetHeader(r *http.Request, name string) string {
	return r.Header.Get(name)
}

func GetIp(r *http.Request) (ip string) {
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
