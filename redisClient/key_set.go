package redisClient

import (
	"strings"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"

	"github.com/garyburd/redigo/redis"
)

// SetSMEMBERS SMEMBERS key
// 返回集合 key 中的所有成员。
// 获取Set value 返回集合 key 中的所有成员。
func (c *RedisClient) SetSMEMBERS(key string) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	log.InfoFTimes(3, "[Redis Log] execute : SMEMBERS %s", key)
	return redis.Values(conn.Do("SMEMBERS", key))
}

// SetSADD SADD 新创建Set  将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
func (c *RedisClient) SetSADD(key string, values []interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	args := redis.Args{}.Add(key)
	for _, value := range values {
		args = args.Add(value)
	}
	log.InfoFTimes(3, "[Redis Log] execute : SADD %s %s", key, strings.Join(utils.AnyToStrings(values), " "))
	_, err = conn.Do("SADD", args...)
	return err
}

// SetSCARD SCARD key
// 返回集合 key 的基数(集合中元素的数量)。
func (c *RedisClient) SetSCARD(key string) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	log.InfoFTimes(3, "[Redis Log] execute : SCARD %s", key)
	_, err = redis.Int64(conn.Do("SCARD ", key))
	return err
}

// SetSDIFF SDIFF key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合之间的差集。
// 不存在的 key 被视为空集。
func (c *RedisClient) SetSDIFF(keys []string) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	args := redis.Args{}
	for _, key := range keys {
		args = args.Add(key)
	}
	log.InfoFTimes(3, "[Redis Log] execute : SDIFF %s", strings.Join(keys, " "))
	return redis.Values(conn.Do("SDIFF", args...))
}

// SetSDIFFSTORE SDIFFSTORE destination key [key ...]
// 这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (c *RedisClient) SetSDIFFSTORE(key string, keys []string) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	args := redis.Args{}.Add(key)
	for _, key := range keys {
		args = args.Add(key)
	}
	log.InfoFTimes(3, "[Redis Log] execute : SDIFFSTORE %s %s", key, strings.Join(keys, " "))
	return redis.Values(conn.Do("SDIFFSTORE", args...))
}

// SetSINTER SINTER key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合的交集。
// 不存在的 key 被视为空集。
func (c *RedisClient) SetSINTER(keys []string) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	args := redis.Args{}
	for _, key := range keys {
		args = args.Add(key)
	}
	log.InfoFTimes(3, "[Redis Log] execute : SINTER %s", strings.Join(keys, " "))
	return redis.Values(conn.Do("SINTER", args...))
}

// SetSINTERSTORE SINTERSTORE destination key [key ...]
// 这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
// 如果 destination 集合已经存在，则将其覆盖。
// destination 可以是 key 本身。
func (c *RedisClient) SetSINTERSTORE(key string, keys []string) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	args := redis.Args{}.Add(key)
	for _, k := range keys {
		args = args.Add(k)
	}
	log.InfoFTimes(3, "[Redis Log] execute : SINTERSTORE %s %s", key, strings.Join(keys, " "))
	return redis.Values(conn.Do("SINTERSTORE", args...))
}

// SetSISMEMBER SISMEMBER key member
// 判断 member 元素是否集合 key 的成员。
// 返回值:
// 如果 member 元素是集合的成员，返回 1 。
// 如果 member 元素不是集合的成员，或 key 不存在，返回 0 。
func (c *RedisClient) SetSISMEMBER(key string, value interface{}) (bool, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return false, err
	}
	log.InfoFTimes(3, "[Redis Log] execute : SISMEMBER %s %v", key, value)
	arg := redis.Args{}.Add(key).Add(value)
	res, err := redis.Int64(conn.Do("SISMEMBER", arg...))
	if err != nil {
		return false, err
	}
	if res == 1 {
		return true, nil
	}
	return false, nil
}

// SetSMOVE SMOVE source destination member
// 将 member 元素从 source 集合移动到 destination 集合。
// SMOVE 是原子性操作。
// 如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。否则，
// member 元素从 source 集合中被移除，并添加到 destination 集合中去。
// 当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
// 当 source 或 destination 不是集合类型时，返回一个错误。
// 返回值: 成功移除，返回 1 。失败0
func (c *RedisClient) SetSMOVE(key, destination string, member interface{}) (bool, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return false, err
	}
	log.InfoFTimes(3, "[Redis Log] execute : SMOVE %s %v %v", key, destination, member)
	arg := redis.Args{}.Add(key).Add(destination).Add(member)
	res, err := redis.Int64(conn.Do("SMOVE", arg...))
	if err != nil {
		return false, err
	}
	if res == 1 {
		return true, nil
	}
	return false, nil
}

// SetSPOP SPOP key
// 移除并返回集合中的一个随机元素。
func (c *RedisClient) SetSPOP(key string) (string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return "", err
	}
	log.InfoFTimes(3, "[Redis Log] execute : SPOP %s", key)
	return redis.String(conn.Do("SPOP", key))
}

// SetSRANDMEMBER SRANDMEMBER key [count]
// 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。
// 如果 count 大于等于集合基数，那么返回整个集合。
// 如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
func (c *RedisClient) SetSRANDMEMBER(key string, count int64) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	arg := redis.Args{}.Add(key).Add(count)
	log.InfoFTimes(3, "[Redis Log] execute : SRANDMEMBER %s %v", key, count)
	return redis.Values(conn.Do("SRANDMEMBER", arg...))
}

// SetSREM SREM key member [member ...]
// 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
func (c *RedisClient) SetSREM(key string, member []interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	args := redis.Args{}.Add(key)
	for _, v := range member {
		args = args.Add(v)
	}
	log.InfoFTimes(3, "[Redis Log] execute : SREM %s %s", key, strings.Join(utils.AnyToStrings(member), " "))
	_, err = conn.Do("SREM", args...)
	return err
}

// SetSUNION SUNION key [key ...]
// 返回一个集合的全部成员，该集合是所有给定集合的并集。
func (c *RedisClient) SetSUNION(keys []string) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	args := redis.Args{}
	for _, v := range keys {
		args = args.Add(v)
	}
	log.InfoFTimes(3, "[Redis Log] execute : SUNION %s", strings.Join(keys, " "))
	return redis.Values(conn.Do("SUNION", args...))
}

// SetSUNIONSTORE SUNIONSTORE destination key [key ...]
// 这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
func (c *RedisClient) SetSUNIONSTORE(key string, keys []string) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	args := redis.Args{}.Add(key)
	for _, v := range keys {
		args = args.Add(v)
	}
	log.InfoFTimes(3, "[Redis Log] execute : SUNIONSTORE %s %s", key, strings.Join(keys, " "))
	return redis.Values(conn.Do("SUNIONSTORE", args...))
}

// 搜索值  SSCAN key cursor [MATCH pattern] [COUNT count]
