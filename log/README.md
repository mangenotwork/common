# 日志

#### 设置日志输出位置
> log.SetLogFile(name string)

- name ：路径+日志名称， 例: ./

---

#### 设置日志输出远程UDP 日志服务
> log.SetOutService(ip string, port int) 

参数
- ip ： 远程ip
- port ： 远程端口

---

#### 全局关闭日志
> log.Close()

---

#### 关闭日志在终端打印
> log.DisableTerminal()

---

#### 输出
> log.Print(args ...interface{})
>
> log.PrintF(format string, args ...interface{})

---

#### 打印信息，日志带有 [INFO]
> log.Info(args ...interface{})
>
> log.InfoF(format string, args ...interface{}) 

> log.InfoTimes(times int, args ...interface{})

参数
- times: runtime.Caller(times)

---

#### 打印Debug，日志带有 [DEBUG]
> log.Debug(args ...interface{})
>
> log.DebugF(format string, args ...interface{})

> log.DebugTimes(times int, args ...interface{})

参数
- times: runtime.Caller(times)

---

#### 打印Warn，日志带有 [WARN] 
> log.Warn(args ...interface{})
>
> log.Warn(format string, args ...interface{})

> log.Warn(times int, args ...interface{})

参数
- times: runtime.Caller(times)

---

#### 打印Error，日志带有 [ERROR] 
> log.Error(args ...interface{})
>
> log.Error(format string, args ...interface{})

> log.Error(times int, args ...interface{})

参数
- times: runtime.Caller(times)

---

#### 打印Panic，日志带有 [PANIC] 

> log.Panic(args ...interface{})


---

