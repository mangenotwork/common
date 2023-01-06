package redisClient

import (
	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"strings"
)

// SetSMEMBERS SMEMBERS key
// 返回集合 key 中的所有成员。
// 获取Set value 返回集合 key 中的所有成员。
func (c *RedisClient) SetSMEMBERS(key string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SMEMBERS ", key)
	return redis.Values(c.Conn.Do("SMEMBERS", key))
}

// SetSADD SADD 新创建Set  将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
func (c *RedisClient) SetSADD(key string, values []interface{}) error {
	if c.Conn == nil {
		return NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, value := range values {
		args = args.Add(value)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SADD ", strings.Join(utils.AnyToStrings(values), " "))
	_, err := c.Conn.Do("SADD", args...)
	return err
}

// SetSCARD SCARD key
// 返回集合 key 的基数(集合中元素的数量)。
func (c *RedisClient) SetSCARD(key string) error {
	if c.Conn == nil {
		return NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SCARD ", key)
	_, err := redis.Int64(c.Conn.Do("SCARD ", key))
	return err
}

// SetSDIFF SDIFF key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合之间的差集。
// 不存在的 key 被视为空集。
func (c *RedisClient) SetSDIFF(keys []string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	args := redis.Args{}
	for _, key := range keys {
		args = args.Add(key)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SDIFF ", strings.Join(keys, " "))
	return redis.Values(c.Conn.Do("SDIFF", args))
}

// SetSDIFFSTORE SDIFFSTORE destination key [key ...]
// 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (c *RedisClient) SetSDIFFSTORE(key string, keys []string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, key := range keys {
		args = args.Add(key)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SDIFFSTORE ", key, " ", strings.Join(keys, " "))
	return redis.Values(c.Conn.Do("SDIFFSTORE", args))
}

// SetSINTER SINTER key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合的交集。
// 不存在的 key 被视为空集。
func (c *RedisClient) SetSINTER(keys []string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	args := redis.Args{}
	for _, key := range keys {
		args = args.Add(key)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SINTER ", strings.Join(keys, " "))
	return redis.Values(c.Conn.Do("SINTER", args))
}

// SetSINTERSTORE SINTERSTORE destination key [key ...]
// 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (c *RedisClient) SetSINTERSTORE(key string, keys []string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, k := range keys {
		args = args.Add(k)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SINTERSTORE ", key, " ", strings.Join(keys, " "))
	return redis.Values(c.Conn.Do("SINTERSTORE", args))
}

// SetSISMEMBER SISMEMBER key member
// 判断 member 元素是否集合 key 的成员。
// 返回值:
// 如果 member 元素是集合的成员，返回 1 。
// 如果 member 元素不是集合的成员，或 key 不存在，返回 0 。
func (c *RedisClient) SetSISMEMBER(key string, value interface{}) (resBool bool, err error) {
	if c.Conn == nil {
		return false, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SISMEMBER ", key, " ", value)
	resBool = false
	res, err := redis.Int64(c.Conn.Do("SISMEMBER", key, value))
	if err != nil {
		return
	}
	if res == 1 {
		resBool = true
		return
	}
	return
}

// SetSMOVE SMOVE source destination member
// 将 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则，
// member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
// 返回值: 成功移除，返回 1 。失败0
func (c *RedisClient) SetSMOVE(key, destination string, member interface{}) (resBool bool, err error) {
	if c.Conn == nil {
		return false, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SMOVE ", key, " ", destination, " ", member)
	resBool = false
	res, err := redis.Int64(c.Conn.Do("SMOVE", key, destination, member))
	if err != nil {
		return
	}
	if res == 1 {
		resBool = true
		return
	}
	return
}

// SetSPOP SPOP key
// 移除并返回集合中的一个随机元素。
func (c *RedisClient) SetSPOP(key string) (string, error) {
	if c.Conn == nil {
		return "", NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SPOP", " ", key)
	return redis.String(c.Conn.Do("SPOP", key))
}

// SetSRANDMEMBER SRANDMEMBER key [count]
// 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。
// 如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
func (c *RedisClient) SetSRANDMEMBER(key string, count int64) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SRANDMEMBER ", key, " ", count)
	return redis.Values(c.Conn.Do("SRANDMEMBER", key, count))
}

// SetSREM SREM key member [member ...]
// 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
func (c *RedisClient) SetSREM(key string, member []interface{}) error {
	if c.Conn == nil {
		return NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range member {
		args = args.Add(v)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SREM ", key,
		strings.Join(utils.AnyToStrings(member), " "))
	_, err := c.Conn.Do("SREM", args)
	return err
}

// SetSUNION SUNION key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合的并集。
func (c *RedisClient) SetSUNION(keys []string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	args := redis.Args{}
	for _, v := range keys {
		args = args.Add(v)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SUNION ", strings.Join(keys, " "))
	return redis.Values(c.Conn.Do("SUNION", args))
}

// SetSUNIONSTORE SUNIONSTORE destination key [key ...]
// 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
func (c *RedisClient) SetSUNIONSTORE(key string, keys []string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range keys {
		args = args.Add(v)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "SUNIONSTORE", key, " ", strings.Join(keys, " "))
	return redis.Values(c.Conn.Do("SUNIONSTORE", args))
}

//搜索值  SSCAN key cursor [MATCH pattern] [COUNT count]
