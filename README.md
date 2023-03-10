# common
常用的方法整理与收集的库，主要目的是为了代码复用，提升工作效率

# 使用文档

- [日志](log/README.md) 
- [验证码](captcha/README.md)
- [配置](conf/README.md)
- [gin相关扩展方法](ginHelper/README.md)
- [http中间件](httpMiddleware/README.md)
- [Mysql](mysqlClient/README.md)
- [Redis](redisClient/README.md)
- [常用](utils/README.md)

# 使用场景
- 与gin配合使用进行开发

# 结构
- httpMiddleware    http中间件
- log       日志打印
- conf      配置的读取
- utils     实用方法
- sshClient       ssh相关的方法，常用与我写本地脚本操作线上服务的相关工作
- mysqlClient     mysql客户端常用方法
- redisClient     redis客户端常用方法
- mq              消息队列常用方法
- osd             对象存储


# 说明
这是一个不断迭代，不断增加方法的一个库，涉及到我个人经历的开发工作中遇到的可复用代码，这个库帮助我提升了大量工作效率，
该库+gin+vue两天写了一个运营操作平台，用于生产当中，目前稳定运行。


# Task
- 使用 golang.org/x/time/rate 做限流中间件
- 对象存储 阿里，腾讯，minio
- 文档生成（excel,pdf,word）
- examples gin+common 基础实例
- 增加使用文档
- 增加使用实例

# TODO
- grpc
- tcp
- udp
- pdf
- http client
- 三方登录（三方授权）
- 三方支付
- 三方短信
- 发送邮件
- 钉钉对接
