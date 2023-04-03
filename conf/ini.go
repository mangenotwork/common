package conf

import (
	"fmt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

// InitConfIni 读取ini文件 获取配置, 常用于 func init() 中
func InitConfIni(path string) {
	confFileName := "app.ini"
	workPath, _ := os.Getwd()
	if os.Getenv("RUNMODE") != "" {
		confFileName = os.Getenv("RUNMODE") + ".ini"
	}
	appConfigPath := filepath.Join(workPath, path, confFileName)
	if !utils.FileExists(appConfigPath) {
		panic("【启动失败】 未找到配置文件!" + appConfigPath)
	}
	log.Print("[启动]读取配置文件:", appConfigPath)
	//读取ini文件到缓存中
	iniObj, err := ini.Load("my.ini")
	if err != nil {
		log.Error("Fail to read file:", err)
		panic("【启动失败】 读取配置文件出现错误!" + appConfigPath)
	}

	Conf.Default.App = &App{
		Name:    iniObj.Section("app").Key("name").String(),
		RunType: iniObj.Section("app").Key("runType").String(),
	}

	httpServerOpen, _ := iniObj.Section("httpServer").Key("open").Bool()
	Conf.Default.HttpServer = &HttpServer{
		Open: httpServerOpen,
		Prod: iniObj.Section("httpServer").Key("prod").String(),
	}

	grpcServerOpen, _ := iniObj.Section("grpcServer").Key("open").Bool()
	grpcServerLog, _ := iniObj.Section("grpcServer").Key("log").Bool()
	Conf.Default.GrpcServer = &GrpcServer{
		Open: grpcServerOpen,
		Prod: iniObj.Section("grpcServer").Key("prod").String(),
		Log:  grpcServerLog,
	}

	Conf.Default.GrpcClient = &GrpcClient{
		Prod: iniObj.Section("grpcClient").Key("prod").String(),
	}

	tcpServerOpen, _ := iniObj.Section("tcpServer").Key("open").Bool()
	Conf.Default.TcpServer = &TcpServer{
		Open: tcpServerOpen,
		Prod: iniObj.Section("tcpServer").Key("prod").String(),
	}

	Conf.Default.TcpClient = &TcpClient{
		Prod: iniObj.Section("tcpClient").Key("prod").String(),
	}

	udpServerOpen, _ := iniObj.Section("udpServer").Key("open").Bool()
	Conf.Default.UdpServer = &UdpServer{
		Open: udpServerOpen,
		Prod: iniObj.Section("udpServer").Key("prod").String(),
	}

	Conf.Default.UdpClient = &UdpClient{
		Prod: iniObj.Section("udpClient").Key("prod").String(),
	}

	redisMaxIdle, _ := iniObj.Section("redis").Key("maxIdle").Int()
	redisMaxActive, _ := iniObj.Section("redis").Key("maxActive").Int()
	Conf.Default.Redis = []*Redis{
		{
			Name:      iniObj.Section("redis").Key("name").String(),
			Host:      iniObj.Section("redis").Key("host").String(),
			Port:      iniObj.Section("redis").Key("port").String(),
			DB:        iniObj.Section("redis").Key("db").String(),
			Password:  iniObj.Section("redis").Key("password").String(),
			MaxIdle:   redisMaxIdle,
			MaxActive: redisMaxActive,
		},
	}

	Conf.Default.Mysql = []*Mysql{
		{
			DBName:   iniObj.Section("mysql").Key("dbname").String(),
			User:     iniObj.Section("mysql").Key("user").String(),
			Password: iniObj.Section("mysql").Key("password").String(),
			Host:     iniObj.Section("mysql").Key("host").String(),
			Port:     iniObj.Section("mysql").Key("port").String(),
		},
	}

	Conf.Default.MqType = iniObj.Section("").Key("mqType").String()

	Conf.Default.Nsq = &Nsq{
		Producer: iniObj.Section("nsq").Key("producer").String(),
		Consumer: iniObj.Section("nsq").Key("consumer").String(),
	}

	Conf.Default.Rabbit = &Rabbit{
		Addr:     iniObj.Section("rabbit").Key("addr").String(),
		User:     iniObj.Section("rabbit").Key("user").String(),
		Password: iniObj.Section("rabbit").Key("password").String(),
	}

	Conf.Default.Kafka = &Kafka{
		Addr: iniObj.Section("kafka").Key("addr").Strings(""),
	}

	Conf.Default.Mongo = []*Mongo{
		{
			Host:     iniObj.Section("mongo").Key("host").String(),
			User:     iniObj.Section("mongo").Key("user").String(),
			Password: iniObj.Section("mongo").Key("password").String(),
		},
	}

	Conf.Default.TTF = iniObj.Section("").Key("ttf").String()

	clusterOpen, _ := iniObj.Section("cluster").Key("open").Bool()
	Conf.Default.Cluster = &Cluster{
		Open:        clusterOpen,
		MyAddr:      iniObj.Section("cluster").Key("myAddr").String(),
		InitCluster: iniObj.Section("cluster").Key("initCluster").String(),
	}

	logCentrePort, _ := iniObj.Section("logCentre").Key("prod").Int()
	Conf.Default.LogCentre = &LogCentre{
		Host: iniObj.Section("logCentre").Key("host").String(),
		Port: logCentrePort,
	}

	JwtExpire, _ := iniObj.Section("jwt").Key("expire").Int()
	Conf.Default.Jwt = &Jwt{
		Secret: iniObj.Section("jwt").Key("secret").String(),
		Expire: JwtExpire,
	}

	Conf.Default.Minio = &Minio{
		Host:   iniObj.Section("minio").Key("host").String(),
		Access: iniObj.Section("minio").Key("access").String(),
		Secret: iniObj.Section("minio").Key("secret").String(),
	}

	Conf.Default.Mq = iniObj.Section("").Key("mq").String()

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

func (c *conf) InitIni() error {
	if !utils.FileExists(c.Path) {
		return fmt.Errorf("未找到配置文件 [%s] !", c.Path)
	}
	log.Info("读取配置文件:", c.Path)
	var err error
	//读取ini文件到缓存中
	c.IniObj, err = ini.Load(c.Path)
	if err != nil {
		log.Error("Fail to read file: ", err)
	}
	return err
}
