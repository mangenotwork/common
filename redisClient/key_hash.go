package redisClient

import (
	"strings"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"

	"github.com/garyburd/redigo/redis"
)

// HashHGETALL HGETALL 获取Hash value
func (c *RedisClient) HashHGETALL(key string) (map[string]string, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoFTimes(3, "[Redis Log] execute : HGETALL %s", key)
	return redis.StringMap(c.Conn.Do("HGETALL", key))
}

// HashHGETALLInt HGETALL 获取Hash value
func (c *RedisClient) HashHGETALLInt(key string) (map[string]int, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoFTimes(3, "[Redis Log] execute : HGETALL %s", key)
	return redis.IntMap(c.Conn.Do("HGETALL", key))
}

// HashHGETALLInt64 HGETALL 获取Hash value
func (c *RedisClient) HashHGETALLInt64(key string) (map[string]int64, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoFTimes(3, "[Redis Log] execute : HGETALL %s", key)
	return redis.Int64Map(c.Conn.Do("HGETALL", key))
}

// HashHSET HSET 新建Hash 单个field
// 如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果域 field 已经存在于哈希表中，旧值将被覆盖。
func (c *RedisClient) HashHSET(key, field string, value interface{}) (int64, error) {
	if c.Conn == nil {
		return 0, NotConnError
	}
	arg := redis.Args{}.Add(key).Add(field).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : HSET %s %v %v", key, field, value)
	return redis.Int64(c.Conn.Do("HSET", arg...))
}

// HashHMSET HMSET 新建Hash 多个field
// HMSET key field value [field value ...]
// 同时将多个 field-value (域-值)对设置到哈希表 key 中。
// 此命令会覆盖哈希表中已存在的域。
func (c *RedisClient) HashHMSET(key string, values map[interface{}]interface{}) error {
	if c.Conn == nil {
		return NotConnError
	}
	args := redis.Args{}.Add(key)
	for k, v := range values {
		args = args.Add(k)
		args = args.Add(v)
	}
	log.InfoFTimes(3, "[Redis Log] execute : HMSET %s %s ", key, strings.Join(utils.AnyToStrings(values), " "))
	_, err := c.Conn.Do("HMSET", args...)
	return err
}

// HashHSETNX HSETNX key field value
// 给hash追加field value
// 将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
func (c *RedisClient) HashHSETNX(key, field string, value interface{}) error {
	if c.Conn == nil {
		return NotConnError
	}
	arg := redis.Args{}.Add(key).Add(field).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : HSETNX %s %v %v", key, field, value)
	_, err := c.Conn.Do("HSETNX", arg...)
	return err
}

// HashHDEL HDEL key field [field ...] 删除哈希表
// key 中的一个或多个指定域，不存在的域将被忽略。
func (c *RedisClient) HashHDEL(key string, fields []string) error {
	if c.Conn == nil {
		return NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range fields {
		args = args.Add(v)
	}
	log.InfoFTimes(3, "[Redis Log] execute : HDEL %s %s", key, strings.Join(fields, " "))
	_, err := c.Conn.Do("HDEL", args...)
	return err
}

// HashHEXISTS HEXISTS key field 查看哈希表
// key 中，给定域 field 是否存在。
func (c *RedisClient) HashHEXISTS(key, fields string) bool {
	if c.Conn == nil {
		log.ErrorTimes(3, "[Redis Log] error :", NotConnError)
		return false
	}
	arg := redis.Args{}.Add(key).Add(fields)
	log.InfoFTimes(3, "[Redis Log] execute : HEXISTS %s %s", key, fields)
	res, err := redis.Int(c.Conn.Do("HEXISTS", arg...))
	if err != nil {
		log.ErrorTimes(3, "[Redis Log] error :", err)
		return false
	}
	if res == 0 {
		return false
	}
	return true
}

// HashHGET HGET key field 返回哈希表
// key 中给定域 field 的值。
func (c *RedisClient) HashHGET(key, fields string) (string, error) {
	if c.Conn == nil {
		return "", NotConnError
	}
	arg := redis.Args{}.Add(key).Add(fields)
	log.InfoFTimes(3, "[Redis Log] execute : HGET %s %s", key, fields)
	return redis.String(c.Conn.Do("HGET", arg...))
}

// HashHINCRBY HINCRBY key field increment
// 为哈希表 key 中的域 field 的值加上增量 increment 。
// 增量也可以为负数，相当于对给定域进行减法操作。
// 如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
// 如果域 field 不存在，那么在执行命令前，域的值被初始化为 0
func (c *RedisClient) HashHINCRBY(key, field string, increment int64) (int64, error) {
	if c.Conn == nil {
		return 0, NotConnError
	}
	arg := redis.Args{}.Add(key).Add(field).Add(increment)
	log.InfoFTimes(3, "[Redis Log] execute : HINCRBY %s %v %v", key, field, increment)
	return redis.Int64(c.Conn.Do("HINCRBY", arg...))
}

// HashHINCRBYFLOAT HINCRBYFLOAT key field increment
// 为哈希表 key 中的域 field 加上浮点数增量 increment 。
// 如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
// 如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
func (c *RedisClient) HashHINCRBYFLOAT(key, field string, increment float64) (float64, error) {
	if c.Conn == nil {
		return 0, NotConnError
	}
	arg := redis.Args{}.Add(key).Add(field).Add(increment)
	log.InfoFTimes(3, "[Redis Log] execute : HINCRBYFLOAT %s %v %v", key, field, increment)
	return redis.Float64(c.Conn.Do("HINCRBYFLOAT", arg...))
}

// HashHKEYS HKEYS key 返回哈希表
// key 中的所有域。
func (c *RedisClient) HashHKEYS(key string) ([]string, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoFTimes(3, "[Redis Log] execute : HKEYS %s", key)
	return redis.Strings(c.Conn.Do("HKEYS", key))
}

// HashHLEN HLEN key 返回哈希表
// key 中域的数量。
func (c *RedisClient) HashHLEN(key string) (int64, error) {
	if c.Conn == nil {
		return 0, NotConnError
	}
	log.InfoFTimes(3, "[Redis Log] execute : HLEN %s", key)
	return redis.Int64(c.Conn.Do("HLEN", key))
}

// HashHMGET HMGET key field [field ...]
// 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
func (c *RedisClient) HashHMGET(key string, fields []string) ([]string, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range fields {
		args = args.Add(v)
	}
	log.InfoFTimes(3, "[Redis Log] execute : HMGET %s %s", key, strings.Join(fields, " "))
	return redis.Strings(c.Conn.Do("HMGET", args...))
}

// HashHVALS HVALS key
// 返回哈希表 key 中所有域的值。
func (c *RedisClient) HashHVALS(key string) ([]string, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoFTimes(3, "[Redis Log] execute : HVALS %s", key)
	return redis.Strings(c.Conn.Do("HVALS", key))
}

// HSCAN
// 搜索value hscan test4 0 match *b*
