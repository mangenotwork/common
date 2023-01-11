# 配置文件

1. 读取配置文件 yaml

2. 读取配置中心 etcd

3. 读取配置中心 自研

4. 读取配置文件 ini

## TODO 文档

## yaml 配置文件
默认读取 app.yaml

环境变量 RUNMODE 确定配置文件昵称

例如： RUNMODE=prod   就会读取prod.yaml文件

```shell
# app相关基础信息
app:
  name: "test_case_1"  # 项目名称
  runType: "dev"  # 运行模式

# http 服务
httpServer:
  open: true
  prod: 12345 # 端口

# grpc 服务
grpcServer:
  open: true
  prod: 12346
  log: true # 打印日志

# grpc 客户端
grpcClient:
  prod: true

# tcp 服务
tcpServer:
  open: true
  prod: 12347

# tcp 客户端
tcpClient:
  prod: 12348

# udp 服务
udpServer:
  open: true
  prod: 12349

# udp 客户端
udpClient:
  prod: 12348

# redis 配置
redis:
  -
    name: "redis1"
    host: "127.0.0.1"
    port: "3306"
    db: 1
    password: ""
    maxIdle: 10  # 最大 Idle 连接
    maxActive: 50 # 最大 活跃 连接

# mysql 配置
mysql:
  -
    dbname: "test"
    user: "root"
    password: "123"
    host: "127.0.0.1"
    port: "3306"


mqType: ""


nsq:
  producer: ""
  consumer: ""

rabbit:
  addr: ""
  user: ""
  password: ""

kafka:
  addr: ""

mongo:
  host: ""
  user: ""
  password: ""

ttf: ""


cluster:
  open: ""
  myAddr: ""
  initCluster: ""


logCentre:
  host: ""
  prod: ""


jwt:
  secret: ""
  expire: ""


minio:
  host: ""
  access: ""
  secret: ""

```

## TODO
- yaml 配置文件 [ok]
- ini 配置文件 
- etcd 配置中心
- 自研配置中心服务， ManGe-ConfigCenter


#### 全局配置对象
> var Config 

---

#### 读取yaml配置文件，将配置保存到全局配置对象 Config
> func NewConf(appConfigPath string) error

---

#### Config 初始化配置对象
> func (c *conf) Init() error

---

#### Config 获取配置输出整形
> func (c *conf) GetInt(key string) int

---

#### Config 获取配置
> func (c *conf) Get(key string) interface{}

---

#### Config 获取配置输出字符串类型
> func (c *conf) GetStr(key string) string

---

#### Configs 默认配置对象, 通用结构, 参考 上面的 ”yaml 配置文件“
```shell
type Configs struct {
	App        *App        `yaml:"app"`
	HttpServer *HttpServer `yaml:"httpServer"`
	GrpcServer *GrpcServer `yaml:"grpcServer"`
	GrpcClient *GrpcClient `yaml:"grpcClient"`
	TcpServer  *TcpServer  `yaml:"tcpServer"`
	TcpClient  *TcpClient  `yaml:"tcpClient"`
	UdpServer  *UdpServer  `yaml:"udpServer"`
	UdpClient  *UdpClient  `yaml:"udpClient"`
	Redis      []*Redis    `yaml:"redis"`
	Mysql      []*Mysql    `yaml:"mysql"`
	MqType     string      `yaml:"mqType"`
	Nsq        *Nsq        `yaml:"nsq"`
	Rabbit     *Rabbit     `yaml:"rabbit"`
	Kafka      *Kafka      `yaml:"kafka"`
	Mongo      []*Mongo    `yaml:"mongo"`
	TTF        string      `yaml:"ttf"`
	Cluster    *Cluster    `yaml:"cluster"`
	LogCentre  *LogCentre  `yaml:"logCentre"`
	Jwt        *Jwt        `yaml:"jwt"`
	Minio      *Minio      `yaml:"minio"`
}
```

---

#### 全局默认配置对象  Conf
> var Conf *Configs = &Configs{}

---

#### 读取yaml文件 获取配置, 常用于 func init() 中
> func InitConf(path string) 

使用
```shell
例如配置文件目录结构
./conf/app.yaml
./conf/dev.yaml
./conf/prod.yaml

InitConf("./conf/")

```

---


