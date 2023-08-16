package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"

	"gopkg.in/yaml.v3"
)

var Conf *Configs = &Configs{}

type Configs struct {
	YamlPath string
	YamlData map[string]interface{}
	Default  *DefaultConf
}

type DefaultConf struct {
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
	Mq         string      `yaml:"mq"`
	User       []*User     `yaml:"user"`
}

// App app相关基础信息
type App struct {
	Name    string `yaml:"name"`
	RunType string `yaml:"runType"` // 项目昵称
}

// HttpServer http服务
type HttpServer struct {
	Open bool   `yaml:"open"`
	Prod string `yaml:"prod"`
}

// GrpcServer grpc服务
type GrpcServer struct {
	Open bool   `yaml:"open"`
	Prod string `yaml:"prod"`
	Log  bool   `yaml:"log"`
}

// GrpcClient grpc客户端
type GrpcClient struct {
	Prod string `yaml:"prod"`
}

// TcpServer tcp服务
type TcpServer struct {
	Open bool   `yaml:"open"`
	Prod string `yaml:"prod"`
}

// TcpClient tcp客户端
type TcpClient struct {
	Prod string `yaml:"prod"`
}

// UdpServer udp服务
type UdpServer struct {
	Open bool   `yaml:"open"`
	Prod string `yaml:"prod"`
}

// UdpClient udp客户端
type UdpClient struct {
	Prod string `yaml:"prod"`
}

// Redis redis配置
type Redis struct {
	Name      string `yaml:"name"` // 自定义一个昵称
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	DB        string `yaml:"db"`
	Password  string `yaml:"password"`
	MaxIdle   int    `yaml:"maxIdle"`
	MaxActive int    `yaml:"maxActive"`
}

// Mysql mysql配置
type Mysql struct {
	DBName   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

// MqType 消息队列类型
type MqType struct {
}

// Nsq 消息队列nsq配置
type Nsq struct {
	Producer string `yaml:"producer"`
	Consumer string `yaml:"consumer"`
}

// Rabbit 消息队列rabbit配置
type Rabbit struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Kafka 消息队列kafka配置
type Kafka struct {
	Addr []string `yaml:"addr"`
}

// Mongo mongo配置
type Mongo struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Cluster 集群使用 主要用于 ServiceTable
type Cluster struct {
	Open        bool   `yaml:"open"`
	MyAddr      string `yaml:"myAddr"`
	InitCluster string `yaml:"initCluster"`
}

// LogCentre 日志服务收集日志配置
type LogCentre struct {
	Host string `yaml:"host"`
	Port int    `yaml:"prod"`
}

// Jwt jwt配置
type Jwt struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"`
}

// Minio 对象存储 minio配置
type Minio struct {
	Host   string `yaml:"host"`
	Access string `yaml:"access"`
	Secret string `yaml:"secret"`
}

// User 默认用户
type User struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

// InitConf 读取yaml文件 获取配置, 常用于 func init() 中
func InitConf(path string) {
	confFileName := "app.yaml"
	workPath, _ := os.Getwd()
	if os.Getenv("RUNMODE") != "" {
		confFileName = os.Getenv("RUNMODE") + ".yaml"
	}
	appConfigPath := filepath.Join(workPath, path, confFileName)
	if !utils.FileExists(appConfigPath) {
		panic("【启动失败】 未找到配置文件!" + appConfigPath)
	}
	log.Print("[启动]读取配置文件:", appConfigPath)
	//读取yaml文件到缓存中
	config, err := ioutil.ReadFile(appConfigPath)
	if err != nil {
		panic("【启动失败】" + err.Error())
	}

	Conf.YamlPath = path
	Conf.YamlData = make(map[string]interface{})
	err = yaml.Unmarshal(config, Conf.YamlData)
	if err != nil {
		panic("【启动失败】" + err.Error())
	}
	Conf.Default = &DefaultConf{}
	err = yaml.Unmarshal(config, Conf.Default)
	if err != nil {
		panic("【启动失败】" + err.Error())
	}
	if Conf.Default.Jwt == nil {
		Conf.Default.Jwt = &Jwt{}
	}
	if Conf.Default.Jwt.Secret == "" {
		Conf.Default.Jwt.Secret = "mange-common"
	}
	if Conf.Default.Jwt.Expire == 0 {
		Conf.Default.Jwt.Expire = 3600 * 24 * 7 // 默认7天
	}
}

func (c *conf) InitYaml() error {
	if !utils.FileExists(c.Path) {
		return fmt.Errorf("未找到配置文件 [%s] !", c.Path)
	}
	log.Info("读取配置文件:", c.Path)
	//读取yaml文件到缓存中
	config, err := ioutil.ReadFile(c.Path)
	if err != nil {
		log.ErrorF("读取配置文件[%s]失败", c.Path)
		return err
	}
	return yaml.Unmarshal(config, c.Data)
}

// YamlGet :: 区分 每一级
func YamlGet(key string) (interface{}, bool) {
	var (
		d  interface{}
		ok bool
	)
	keyList := strings.Split(key, "::")
	temp := make(map[string]interface{})
	temp = Conf.YamlData
	for _, v := range keyList {
		d, ok = temp[v]
		if !ok {
			break
		}
		temp = utils.AnyToMap(d)
	}
	return d, ok
}

// YamlGetString  区分
func YamlGetString(key string) (string, error) {
	data, ok := YamlGet(key)
	if ok {
		return utils.AnyToString(data), nil
	}
	return "", fmt.Errorf("配置文件没有找到参数 %s", key)
}

func YamlGetInt64(key string) (int64, error) {
	data, ok := YamlGet(key)
	if ok {
		return utils.AnyToInt64(data), nil
	}
	return 0, fmt.Errorf("配置文件没有找到参数 %s", key)
}
