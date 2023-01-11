# gin 相关的方法，接口

#### 统一接口输出
```shell
// ResponseJson 统一接口输出
type ResponseJson struct {
	Code      int64       `json:"code"`
	Msg       string      `json:"msg"`
	Date      interface{} `json:"data"`
	TimeStamp int64       `json:"timeStamp"`
}
```

---

#### 统一接口输出方法
> func APIOutPut(c *gin.Context, code int64, msg string, data interface{})

---

#### 统一接口输出错误方法
> func APIOutPutError(c *gin.Context, code int64, msg string)

---

#### 获取参数
> func GetPostArgs(c *gin.Context, obj interface{}) error 

---
