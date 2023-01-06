package redisClient

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"golang.org/x/crypto/ssh"
	"net"
)

var NotConnError = fmt.Errorf("未连接redis")

// RedisConnMap 全局redis 连接对象 map
var RedisConnMap map[string]*RedisClient = make(map[string]*RedisClient, 0)

// Client Conn 全局redis 连接对象
var Client *RedisClient = &RedisClient{}

type RedisClient struct {
	Conn redis.Conn
}

func GetClient(name string) (*RedisClient, error) {
	c, ok := RedisConnMap[name]
	if !ok {
		return nil, fmt.Errorf("没有这个redis连接")
	}
	return c, nil
}

// InitConfRedisConn 配置文件读redis配置，并将连接保存到 RedisConnMap
func InitConfRedisConn() error {
	if conf.Conf == nil || len(conf.Conf.Redis) < 1 {
		return fmt.Errorf("未读取到redis配置")
	}
	for _, v := range conf.Conf.Redis {
		conn, err := Conn(v.Host, v.Password, v.DB)
		if err != nil {
			log.Error(err)
			continue
		}
		RedisConnMap[v.Name] = &RedisClient{
			Conn: conn,
		}
	}
	return nil
}

func Conn(host, password, db string) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", host)
	if nil != err {
		return nil, err
	}
	if _, authErr := conn.Do("AUTH", password); authErr != nil {
		return nil, fmt.Errorf("redis auth password error: %s", authErr)
	}
	_, err = conn.Do("select", db)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func RedisConn(host, password, db string) error {
	var err error
	Client.Conn, err = redis.Dial("tcp", host)
	if nil != err {
		return err
	}
	if _, authErr := Client.Conn.Do("AUTH", password); authErr != nil {
		return fmt.Errorf("redis auth password error: %s", authErr)
	}
	if db == "" {
		db = "0"
	}
	return Client.SelectDB(db)
}

// SelectDB 指定db的连接
func (c *RedisClient) SelectDB(db string) error {
	_, err := c.Conn.Do("select", db)
	return err
}

func getSSHClient(user string, pass string, addr string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshConn, err := net.Dial("tcp", addr)
	if nil != err {
		return nil, err
	}
	clientConn, chans, reqs, err := ssh.NewClientConn(sshConn, addr, config)
	if nil != err {
		sshConn.Close()
		return nil, err
	}
	client := ssh.NewClient(clientConn, chans, reqs)
	return client, nil
}

func RedisSshConn(user, pass, addr, host, password, db string) error {
	var err error
	sshClient, err := getSSHClient(user, pass, addr)
	if nil != err {
		return err
	}
	conn, err := sshClient.Dial("tcp", host)
	if nil != err {
		return err
	}
	Client.Conn = redis.NewConn(conn, -1, -1)
	if _, authErr := Client.Conn.Do("AUTH", password); authErr != nil {
		return fmt.Errorf("redis auth password error: %s", authErr)
	}
	if db == "" {
		db = "0"
	}
	return Client.SelectDB(db)
}
