package redisClient

import (
	"strings"

	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"

	"github.com/garyburd/redigo/redis"
)

// ListLRANGEALL LRANGE 获取List value
func (c *RedisClient) ListLRANGEALL(key string) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	log.InfoFTimes(3, "[Redis Log] execute : LRANGE %s 0 -1", key)
	return redis.Values(conn.Do("LRANGE", key, 0, -1))
}

// ListLRANGE LRANGE key start stop
// 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
func (c *RedisClient) ListLRANGE(key string, start, stop int64) ([]interface{}, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	arg := redis.Args{}.Add(key).Add(start).Add(stop)
	log.InfoFTimes(3, "[Redis Log] execute : LRANGE %s %v %v", key, start, stop)
	return redis.Values(conn.Do("LRANGE", arg...))
}

// ListLPUSH LPUSH 新创建list 将一个或多个值 value 插入到列表 key 的表头
func (c *RedisClient) ListLPUSH(key string, values []interface{}) error {
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
	log.InfoFTimes(3, "[Redis Log] execute : LPUSH %s %s", key, strings.Join(utils.AnyToStrings(values), " "))
	_, err = conn.Do("LPUSH", args...)
	return err
}

// ListRPUSH RPUSH key value [value ...]
// 将一个或多个值 value 插入到列表 key 的表尾(最右边)。
// 如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行
// RPUSH mylist a b c ，得出的结果列表为 a b c ，等同于执行命令 RPUSH mylist a 、 RPUSH mylist b 、 RPUSH mylist c 。
// 新创建List  将一个或多个值 value 插入到列表 key 的表尾(最右边)。
func (c *RedisClient) ListRPUSH(key string, values []interface{}) error {
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
	log.InfoFTimes(3, "[Redis Log] execute : RPUSH %s %s", key, strings.Join(utils.AnyToStrings(values), " "))
	_, err = conn.Do("RPUSH", args...)
	return err
}

// ListBLPOP BLPOP key [key ...] timeout
// BLPOP 是列表的阻塞式(blocking)弹出原语。
// 它是 LPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BLPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
func (c *RedisClient) ListBLPOP() {}

// ListBRPOP BRPOP key [key ...] timeout
// BRPOP 是列表的阻塞式(blocking)弹出原语。
// 它是 RPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BRPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
func (c *RedisClient) ListBRPOP() {}

// ListBRPOPLPUSH BRPOPLPUSH source destination timeout
// BRPOPLPUSH 是 RPOPLPUSH 的阻塞版本，当给定列表 source 不为空时， BRPOPLPUSH 的表现和 RPOPLPUSH 一样。
// 当列表 source 为空时， BRPOPLPUSH 命令将阻塞连接，直到等待超时，或有另一个客户端对 source 执行 LPUSH 或 RPUSH 命令为止。
func (c *RedisClient) ListBRPOPLPUSH() {}

// ListLINDEX LINDEX key index
// 返回列表 key 中，下标为 index 的元素。
func (c *RedisClient) ListLINDEX(key string, index int64) (string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return "", err
	}
	arg := redis.Args{}.Add(key).Add(index)
	log.InfoFTimes(3, "[Redis Log] execute : LINDEX %s %v", key, index)
	return redis.String(conn.Do("LINDEX", arg...))
}

// ListLINSERT LINSERT key BEFORE|AFTER pivot value
// 将值 value 插入到列表 key 当中，位于值 pivot 之前或之后。
// 当 pivot 不存在于列表 key 时，不执行任何操作。
// 当 key 不存在时， key 被视为空列表，不执行任何操作。
// 如果 key 不是列表类型，返回一个错误。
// direction : 方向 bool true:BEFORE(前)    false: AFTER(后)
func (c *RedisClient) ListLINSERT(direction bool, key, pivot, value string) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	directionStr := "AFTER"
	if direction {
		directionStr = "BEFORE"
	}
	arg := redis.Args{}.Add(key).Add(directionStr).Add(pivot).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : LINSERT %s %v %v %v", key, directionStr, pivot, value)
	_, err = conn.Do("LINSERT", arg...)
	return err
}

// ListLLEN LLEN key
// 返回列表 key 的长度。
// 如果 key 不存在，则 key 被解释为一个空列表，返回 0 .
func (c *RedisClient) ListLLEN(key string) (int64, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0, err
	}
	log.InfoFTimes(3, "[Redis Log] execute : LLEN %s", key)
	return redis.Int64(conn.Do("LLEN", key))
}

// ListLPOP LPOP key
// 移除并返回列表 key 的头元素。
func (c *RedisClient) ListLPOP(key string) (string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return "", err
	}
	log.InfoFTimes(3, "[Redis Log] execute : LPOP %s", key)
	return redis.String(conn.Do("LPOP", key))
}

// ListLPUSHX LPUSHX key value
// 将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。
// 和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。
func (c *RedisClient) ListLPUSHX(key string, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : LPUSHX %s %v", key, value)
	_, err = conn.Do("LPUSHX", arg...)
	return err
}

// ListLREM LREM key count value
// 根据参数 count 的值，移除列表中与参数 value 相等的元素。
// count 的值可以是以下几种：
// count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。
// count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。
// count = 0 : 移除表中所有与 value 相等的值。
func (c *RedisClient) ListLREM(key string, count int64, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(count).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : LREM %s %v %v", key, count, value)
	_, err = conn.Do("LREM", arg...)
	return err
}

// ListLSET LSET key index value
// 将列表 key 下标为 index 的元素的值设置为 value 。
// 当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。
func (c *RedisClient) ListLSET(key string, index int64, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(index).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : LSET %s %v %v", key, index, value)
	_, err = conn.Do("LSET", arg...)
	return err
}

// ListLTRIM LTRIM key start stop
// 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
// 举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
func (c *RedisClient) ListLTRIM(key string, start, stop int64) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(start).Add(stop)
	log.InfoFTimes(3, "[Redis Log] execute : LTRIM %s %v %v", key, start, stop)
	_, err = conn.Do("LTRIM", arg...)
	return err
}

// ListRPOP RPOP key
// 移除并返回列表 key 的尾元素。
func (c *RedisClient) ListRPOP(key string) (string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return "", err
	}
	log.InfoFTimes(3, "[Redis Log] execute : RPOP %s", key)
	return redis.String(conn.Do("RPOP", key))
}

// ListRPOPLPUSH RPOPLPUSH source destination
// 命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：
// 将列表 source 中的最后一个元素(尾元素)弹出，并返回给客户端。
// 将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。
// 举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination
// 列表有元素 x, y, z ，执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ，
// destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。
// 如果 source 不存在，值 nil 被返回，并且不执行其他动作。
// 如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转(rotation)操作。
func (c *RedisClient) ListRPOPLPUSH(key, destination string) (string, error) {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return "", err
	}
	arg := redis.Args{}.Add(key).Add(destination)
	log.InfoFTimes(3, "[Redis Log] execute : RPOPLPUSH %s %v", key, destination)
	return redis.String(conn.Do("RPOPLPUSH", arg...))
}

// ListRPUSHX RPUSHX key value
// 将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。
func (c *RedisClient) ListRPUSHX(key string, value interface{}) error {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return err
	}
	arg := redis.Args{}.Add(key).Add(value)
	log.InfoFTimes(3, "[Redis Log] execute : RPUSHX %s %v", key, value)
	_, err = conn.Do("RPUSHX", arg...)
	return err
}
