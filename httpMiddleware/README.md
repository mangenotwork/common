# http中间件


#### 获取ip
> httpMiddleware.GetIP(r *http.Request) (ip string)

---

####  CorsHandler 跨域中间件
> httpMiddleware.CorsHandler()

---

#### 启动多个web服务进程，做本地负载
> httpMiddleware.SetSockOptInt() net.ListenConfig

```go
func RunHttpService() {
	var lc = httpMiddleware.SetSockOptInt()
	// 启动五个进程
	for i := 0; i < 5; i++ {
		go func(i int) {
			gin.SetMode(gin.ReleaseMode)
			s := Routers(fmt.Sprintf("%d", i))
			lis, err := lc.Listen(context.Background(), "tcp", "0.0.0.0:14444")
			if err != nil {
				panic("启动 http api 失败, err =  " + err.Error())
			}
			logger.Info("启动 Http API , ID:", i)
			err = s.RunListener(lis)
			if err != nil {
				panic("启动 http api 失败, err =  " + err.Error())
			}
		}(i)
	}
}
```

---