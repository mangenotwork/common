package redisClient

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
)

var NotConnError = fmt.Errorf("未连接redis")

// RedisPoolMap 全局redis 连接对象 map
var RedisPoolMap map[string]*RedisPool = make(map[string]*RedisPool, 0)

// ClientPool 全局redis 连接池
var ClientPool *RedisPool = &RedisPool{}

type RedisClient struct {
	Name string
}

func NewRedisClient(name string) *RedisClient {
	return &RedisClient{
		Name: name,
	}
}

type RedisPool struct {
	Pool *redis.Pool
}

func GetConn(name string) (redis.Conn, error) {
	pool, ok := RedisPoolMap[name]
	if !ok {
		return nil, fmt.Errorf("没有这个redis连接")
	}
	return pool.Pool.Get(), nil
}

func InitConfRedisPool() {
	if conf.Conf == nil || len(conf.Conf.Default.Redis) < 1 {
		panic("未读取到redis配置")
	}
	for _, v := range conf.Conf.Default.Redis {
		conn := RedisPoolConn(fmt.Sprintf("%s:%s", v.Host, v.Port), v.Password, v.DB, v.MaxIdle, v.MaxActive)
		RedisPoolMap[v.Name] = &RedisPool{Pool: conn}
	}
	log.Print("连接Redis : ", RedisPoolMap)
}

// SelectDB 指定db的连接
func (c *RedisClient) SelectDB(db string) error {
	conn, err := GetConn(c.Name)
	_, err = conn.Do("select", db)
	return err
}
