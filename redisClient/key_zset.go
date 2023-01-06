package redisClient

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"strings"
)

// ZSetZRANGEALL ZRANGE 获取ZSet value 返回集合 有序集成员的列表。
func (c *RedisClient) ZSetZRANGEALL(key string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZRANGE ", key, " ", 0, " ", -1, " ", "WITHSCORES")
	return redis.Values(c.Conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
}

// ZSetZRANGE ZRANGE key start stop [WITHSCORES]
// 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
func (c *RedisClient) ZSetZRANGE(key string, start, stop int64) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZRANGE ", key, " ", start, " ", stop, " ", "WITHSCORES")
	return redis.Values(c.Conn.Do("ZRANGE", key, start, stop, "WITHSCORES"))
}

// ZSetZREVRANGE ZREVRANGE key start stop [WITHSCORES]
// 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
func (c *RedisClient) ZSetZREVRANGE(key string, start, stop int64) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZREVRANGE", key, " ", start, " ", stop, " ", "WITHSCORES")
	return redis.Values(c.Conn.Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
}

// ZSetZADD ZADD 新创建ZSet 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
func (c *RedisClient) ZSetZADD(key string, values []interface{}) error {
	if c.Conn == nil {
		return NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, value := range values {
		fmt.Println(value)
		for k, v := range value.(map[string]interface{}) {
			args = args.Add(v)
			args = args.Add(k)
		}
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZADD", key, " ",
		strings.Join(utils.AnyToStrings(args), " "))
	_, err := c.Conn.Do("ZADD", args...)
	return err
}

// ZSetZCARD ZCARD key
// 返回有序集 key 的基数。
func (c *RedisClient) ZSetZCARD(key string) (int64, error) {
	if c.Conn == nil {
		return 0, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZCARD", key)
	return redis.Int64(c.Conn.Do("ZCARD", key))
}

// ZSetZCOUNT ZCOUNT key min max
// 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (c *RedisClient) ZSetZCOUNT(key string, min, max int64) (int64, error) {
	if c.Conn == nil {
		return 0, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZCOUNT ", key, " ", min, " ", max)
	return redis.Int64(c.Conn.Do("ZCOUNT", key, min, max))
}

// ZSetZINCRBY ZINCRBY key increment member
// 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5 。
// 当 key 不存在，或 member 不是 key 的成员时， ZINCRBY key increment member 等同于 ZADD key increment member 。
func (c *RedisClient) ZSetZINCRBY(key, member string, increment int64) (string, error) {
	if c.Conn == nil {
		return "", NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZINCRBY", key, increment, member)
	return redis.String(c.Conn.Do("ZINCRBY", key, increment, member))
}

// ZSetZRANGEBYSCORE ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
// 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
// 具有相同 score 值的成员按字典序(lexicographical order)来排列(该属性是有序集提供的，不需要额外的计算)。
// 可选的 LIMIT 参数指定返回结果的数量及区间(就像SQL中的 SELECT LIMIT offset, count )，注意当 offset 很大时，
// 定位 offset 的操作可能需要遍历整个有序集，此过程最坏复杂度为 O(N) 时间。
// 可选的 WITHSCORES 参数决定结果集是单单返回有序集的成员，还是将有序集成员及其 score 值一起返回。
// 区间及无限
// min 和 max 可以是 -inf 和 +inf ，这样一来，你就可以在不知道有序集的最低和最高 score 值的情况下，使用 ZRANGEBYSCORE 这类命令。
// 默认情况下，区间的取值使用闭区间 (小于等于或大于等于)，你也可以通过给参数前增加 ( 符号来使用可选的开区间 (小于或大于)。
func (c *RedisClient) ZSetZRANGEBYSCORE(key string, min, max, offset, count int64) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZRANGEBYSCORE ", key, " ", min, " ", max, " ",
		offset, " ", count)
	return redis.Values(c.Conn.Do("ZRANGEBYSCORE", key, min, max, offset, count))
}

func (c *RedisClient) ZSetZRANGEBYSCOREALL(key string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZRANGEBYSCORE ", key, " -inf ", "+inf")
	return redis.Values(c.Conn.Do("ZRANGEBYSCORE", key, "-inf", "+inf"))
}

// ZREVRANGEBYSCORE key max min [WITHSCORES] [LIMIT offset count]
// 返回有序集 key 中， score 值介于 max 和 min 之间(默认包括等于 max 或 min )的所有的成员。有序集成员按 score 值递减(从大到小)的次序排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order )排列。

func (c *RedisClient) ZSetZREVRANGEBYSCORE(key string, min, max, offset, count int64) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZREVRANGEBYSCORE ", key, " ", min, " ", max, " ",
		offset, " ", count)
	return redis.Values(c.Conn.Do("ZREVRANGEBYSCORE", key, min, max, offset, count))
}

func (c *RedisClient) ZSetZREVRANGEBYSCOREALL(key string) ([]interface{}, error) {
	if c.Conn == nil {
		return nil, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZREVRANGEBYSCORE ", key, " -inf ", "+inf")
	return redis.Values(c.Conn.Do("ZREVRANGEBYSCORE", key, "-inf", "+inf"))
}

// ZSetZRANK ZRANK key member
// 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
func (c *RedisClient) ZSetZRANK(key string, member interface{}) (int64, error) {
	if c.Conn == nil {
		return 0, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZRANK", key, member)
	return redis.Int64(c.Conn.Do("ZRANK", key, member))
}

// ZSetZREM ZREM key member [member ...]
// 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
func (c *RedisClient) ZSetZREM(key string, member []interface{}) error {
	if c.Conn == nil {
		return NotConnError
	}
	args := redis.Args{}.Add(key)
	for _, v := range member {
		args = args.Add(v)
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZREM ", key,
		strings.Join(utils.AnyToStrings(member), " "))
	_, err := c.Conn.Do("ZREM", args...)
	return err
}

// ZSetZREMRANGEBYRANK ZREMRANGEBYRANK key start stop
// 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
func (c *RedisClient) ZSetZREMRANGEBYRANK(key string, start, stop int64) error {
	if c.Conn == nil {
		return NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZREMRANGEBYRANK ", key, " ", start, " ", stop)
	_, err := redis.Int64(c.Conn.Do("ZREMRANGEBYRANK", key, start, stop))
	return err
}

// ZSetZREMRANGEBYSCORE ZREMRANGEBYSCORE key min max
// 移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
func (c *RedisClient) ZSetZREMRANGEBYSCORE(key string, min, max int64) error {
	if c.Conn == nil {
		return NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZREMRANGEBYSCORE ", key, " ", min, " ", max)
	_, err := c.Conn.Do("ZREMRANGEBYSCORE", key, min, max)
	return err
}

// ZSetZREVRANK ZREVRANK key member
// 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
// 使用 ZRANK 命令可以获得成员按 score 值递增(从小到大)排列的排名。
func (c *RedisClient) ZSetZREVRANK(key string, member interface{}) (int64, error) {
	if c.Conn == nil {
		return 0, NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZREVRANK ", key, " ", member)
	return redis.Int64(c.Conn.Do("ZREVRANK", key, member))
}

// ZSetZSCORE ZSCORE key member
// 返回有序集 key 中，成员 member 的 score
func (c *RedisClient) ZSetZSCORE(key string, member interface{}) (string, error) {
	if c.Conn == nil {
		return "", NotConnError
	}
	log.InfoTimes(3, "[Redis Log] execute :", "ZSCORE ", key, " ", member)
	return redis.String(c.Conn.Do("ZSCORE", key, member))
}

// ZSetZUNIONSTORE ZUNIONSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]
// 计算给定的一个或多个有序集的并集，其中给定 key 的数量必须以 numkeys 参数指定，并将该并集(结果集)储存到 destination 。
func (c *RedisClient) ZSetZUNIONSTORE() {}

// ZSetZINTERSTORE ZINTERSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]
// 计算给定的一个或多个有序集的交集，其中给定 key 的数量必须以 numkeys 参数指定，并将该交集(结果集)储存到 destination 。
// 默认情况下，结果集中某个成员的 score 值是所有给定集下该成员 score 值之和.
func (c *RedisClient) ZSetZINTERSTORE() {}

//搜索值  ZSCAN key cursor [MATCH pattern] [COUNT count]
