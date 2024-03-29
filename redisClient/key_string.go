package redisClient

import (
	"strings"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"

	"github.com/garyburd/redigo/redis"
)

// StringGet GET 获取String value
func (c *RedisClient) StringGet(key string) (string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return "", err
	}
	log.InfoFTimes(3, "[Redis Log] execute : GET %s", key)
	return redis.String(conn.Do("GET", key))
}

// StringSET SET 新建String
func (c *RedisClient) StringSET(key string, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : SET %s %v", key, value)
	_, err = conn.Do("SET", arg...)
	return err
}

// StringSETEX SETEX 新建String 含有时间
func (c *RedisClient) StringSETEX(key string, ttl int64, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(ttl).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : SETEX  %s %v %v", key, ttl, value)
	_, err = conn.Do("SETEX", arg...)
	return err
}

// StringPSETEX PSETEX key milliseconds value
// 这个命令和 SETEX 命令相似，但它以毫秒为单位设置 key 的生存时间，而不是像 SETEX 命令那样，以秒为单位。
func (c *RedisClient) StringPSETEX(key string, ttl int64, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(ttl).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : PSETEX %s %v %v", key, ttl, value)
	_, err = conn.Do("PSETEX", arg...)
	return err
}

// StringSETNX key value
// 将 key 的值设为 value ，当且仅当 key 不存在。
// 若给定的 key 已经存在，则 SETNX 不做任何动作。
func (c *RedisClient) StringSETNX(key string, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	log.InfoFTimes(3, "[Redis Log] execute : SETNX %s %v", key, value)
	arg := redis.Args{}.Add(key).Add(value)
	_, err = conn.Do("SETNX", arg...)
	return err
}

// StringSETRANGE SETRANGE key offset value
// 用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。
// 不存在的 key 当作空白字符串处理。
func (c *RedisClient) StringSETRANGE(key string, offset int64, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(offset).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : SETRANGE %s %v %v", key, offset, value)
	_, err = conn.Do("SETRANGE", arg...)
	return err
}

// StringAPPEND APPEND key value
// 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func (c *RedisClient) StringAPPEND(key string, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : APPEND %s %v", key, value)
	_, err = redis.String(conn.Do("APPEND", arg...))
	return err
}

// StringSETBIT SETBIT key offset value
// 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
// value : 位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
// 注意 offset 不能太大，越大key越大
func (c *RedisClient) StringSETBIT(key string, offset, value int64) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(offset).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : SETBIT %s %d %d", key, offset, value)
	_, err = conn.Do("SETBIT", arg...)
	return err
}

// StringBITCOUNT BITCOUNT key [start] [end]
// 计算给定字符串中，被设置为 1 的比特位的数量。
func (c *RedisClient) StringBITCOUNT(key string) (int64, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0, err
	}
	log.InfoFTimes(3, "[Redis Log] execute : BITCOUNT %s", key)
	return redis.Int64(conn.Do("BITCOUNT", key))
}

// StringGETBIT GETBIT key offset
// 对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
// 当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
func (c *RedisClient) StringGETBIT(key string, offset int64) (int64, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0, err
	}
	arg := redis.Args{}.Add(key).Add(offset)
	log.InfoFTimes(3, "[Redis Log] execute : GETBIT %s %d", key, offset)
	return redis.Int64(conn.Do("GETBIT", arg...))
}

// StringBITOP BITOP operation destkey key [key ...]
// 对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destkey 上。
// BITOP AND destkey key [key ...] ，对一个或多个 key 求逻辑并，并将结果保存到 destkey 。
// BITOP OR destkey key [key ...] ，对一个或多个 key 求逻辑或，并将结果保存到 destkey 。
// BITOP XOR destkey key [key ...] ，对一个或多个 key 求逻辑异或，并将结果保存到 destkey 。
// BITOP NOT destkey key ，对给定 key 求逻辑非，并将结果保存到 destkey 。
func (c *RedisClient) StringBITOP() {}

// StringDECR key
// 将 key 中储存的数字值减一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
func (c *RedisClient) StringDECR(key string) (int64, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0, err
	}
	log.InfoFTimes(3, "[Redis Log] execute : DECR %s", key)
	return redis.Int64(conn.Do("DECR", key))
}

// StringDECRBY DECRBY key decrement
// 将 key 所储存的值减去减量 decrement 。
func (c *RedisClient) StringDECRBY(key, decrement string) (int64, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0, err
	}
	arg := redis.Args{}.Add(key).Add(decrement)
	log.InfoFTimes(3, "[Redis Log] execute : DECRBY %s %v", key, decrement)
	return redis.Int64(conn.Do("DECRBY", arg...))
}

// StringGETRANGE GETRANGE key start end
// 返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
func (c *RedisClient) StringGETRANGE(key string, start, end int64) (string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return "", err
	}
	arg := redis.Args{}.Add(key).Add(start).Add(end)
	log.InfoFTimes(3, "[Redis Log] execute : GETRANGE %s %v %v", key, start, end)
	return redis.String(conn.Do("GETRANGE", arg...))
}

// StringGETSET GETSET key value
// 将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
// 当 key 存在但不是字符串类型时，返回一个错误。
func (c *RedisClient) StringGETSET(key string, value interface{}) (string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return "", err
	}
	arg := redis.Args{}.Add(key).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : GETSET %s %v", key, value)
	return redis.String(conn.Do("GETSET", arg...))
}

// StringINCR INCR key
// 将 key 中储存的数字值增一。
// 如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
// 如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
func (c *RedisClient) StringINCR(key string) (int64, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0, err
	}
	log.InfoFTimes(3, "[Redis Log] execute : INCR %s", key)
	return redis.Int64(conn.Do("INCR", key))
}

// StringINCRBY INCRBY key increment
// 将 key 所储存的值加上增量 increment 。
func (c *RedisClient) StringINCRBY(key, increment string) (int64, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0, err
	}
	arg := redis.Args{}.Add(key).Add(increment)
	log.InfoFTimes(3, "[Redis Log] execute : INCRBY %s %v", key, increment)
	return redis.Int64(conn.Do("INCRBY", arg...))
}

// StringINCRBYFLOAT INCRBYFLOAT key increment
// 为 key 中所储存的值加上浮点数增量 increment 。
func (c *RedisClient) StringINCRBYFLOAT(key, increment float64) (float64, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0, err
	}
	arg := redis.Args{}.Add(key).Add(increment)
	log.InfoFTimes(3, "[Redis Log] execute : INCRBYFLOAT %s %v", key, increment)
	return redis.Float64(conn.Do("INCRBYFLOAT", arg...))
}

// StringMGET MGET key [key ...]
// 返回所有(一个或多个)给定 key 的值。
// 如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。因此，该命令永不失败。
func (c *RedisClient) StringMGET(key []interface{}) ([]string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	args := redis.Args{}
	for _, value := range key {
		args = args.Add(value)
	}
	log.InfoFTimes(3, "[Redis Log] execute : MGET %s %s", key, strings.Join(utils.AnyToStrings(key), " "))
	return redis.Strings(conn.Do("MGET", args...))
}

// StringMSET MSET key value [key value ...]
// 同时设置一个或多个 key-value 对。
// 如果某个给定 key 已经存在，那么 MSET 会用新值覆盖原来的旧值，如果这不是你所希望的效果，
// 请考虑使用 MSETNX 命令：它只会在所有给定 key 都不存在的情况下进行设置操作。
// MSET 是一个原子性(atomic)操作，所有给定 key 都会在同一时间内被设置，某些给定 key 被更新而另一些给定 key 没有改变的情况，不可能发生。
func (c *RedisClient) StringMSET(values []interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	args := redis.Args{}
	for _, value := range values {
		args = args.Add(value)
	}
	log.InfoFTimes(3, "[Redis Log] execute : MSET %s", strings.Join(utils.AnyToStrings(values), " "))
	_, err = conn.Do("MSET", args...)
	return err
}

// StringMSETNX MSETNX key value [key value ...]
// 同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
// 即使只有一个给定 key 已存在， MSETNX 也会拒绝执行所有给定 key 的设置操作。
// MSETNX 是原子性的，因此它可以用作设置多个不同 key 表示不同字段(field)的唯一性逻辑对象(unique logic object)，
// 所有字段要么全被设置，要么全不被设置。
func (c *RedisClient) StringMSETNX(values []interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	args := redis.Args{}
	for _, value := range values {
		args = args.Add(value)
	}
	log.InfoFTimes(3, "[Redis Log] execute : MSETNX %s", strings.Join(utils.AnyToStrings(values), " "))
	_, err = conn.Do("MSETNX", args...)
	return err
}

// StringSTRLEN STRLEN key
// 返回 key 所储存的字符串值的长度。
// 当 key 储存的不是字符串值时，返回一个错误。
func (c *RedisClient) StringSTRLEN() {}
