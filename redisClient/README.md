# redis

#### RedisConnMap 全局 redis连接map
> var RedisConnMap map[string]*RedisClient

---

#### Client 全局 redis连接
> var Client *RedisClient

---

#### RedisClient redis客户端对象
```
type RedisClient struct {
	Conn redis.Conn
}
```

---

#### 获取连接，取的是 RedisConnMap里的redis客户端对象
> func GetClient(name string) (*RedisClient, error)

---

#### 配置文件读redis配置，并将连接保存到 RedisConnMap
> func InitConfRedisConn() error 

配置
```yaml
# redis 配置
redis:
  -
    name: "redis1"
    host: "127.0.0.1"
    port: "3306"
    db: 1
    password: ""
    maxIdle: 10  # 最大 Idle 连接
    maxActive: 50 # 最大 活跃 连接
  -
    name: "redis2"
    host: "127.0.0.1"
    port: "3306"
    db: 2
    password: ""
    maxIdle: 10  # 最大 Idle 连接
    maxActive: 50 # 最大 活跃 连接
```

---

#### 连接redis
> func Conn(host, password, db string) (redis.Conn, error)

---

#### 连接redis, 连接存储到全局redis连接 Client
> func RedisConn(host, password, db string) error

---


#### 指定db的连接
> func (c *RedisClient) SelectDB(db string) error

---

#### 通过ssh隧道连接redis
> func RedisSshConn(user, pass, addr, host, password, db string) error

---

#### GetKeyType 获取key的类型
> func (c *RedisClient) GetKeyType(key string) string

---

#### GetKeyTTL 获取key的过期时间
> func (c *RedisClient) GetKeyTTL(key string) int64 

---

#### EXISTSKey 检查给定 key 是否存在
> func (c *RedisClient) EXISTSKey(key string) bool

---

#### RenameKey 修改key名称
> func (c *RedisClient) RenameKey(name, newName string) bool

--- 

#### UpdateKeyTTL 更新key ttl
> func (c *RedisClient) UpdateKeyTTL(key string, ttl int64) bool

--- 

#### EXPIREATKey 指定key多久过期 接收的是unix时间戳
> func (c *RedisClient) EXPIREATKey(key string, date int64) bool

---

#### DELKey 删除key
> func (c *RedisClient) DELKey(key string) bool

---

#### HashHGETALL HGETALL 获取Hash value
> func (c *RedisClient) HashHGETALL(key string) (map[string]string, error)

---

#### HashHGETALLInt HGETALL 获取Hash value
> func (c *RedisClient) HashHGETALLInt(key string) (map[string]int, error)

----

#### HashHGETALLInt64 HGETALL 获取Hash value
> func (c *RedisClient) HashHGETALLInt64(key string) (map[string]int64, error)

----

#### HashHSET HSET 新建Hash 单个field
> func (c *RedisClient) HashHSET(key, field string, value interface{}) (int64, error)

----

#### HashHMSET HMSET 新建Hash 多个field
> func (c *RedisClient) HashHMSET(key string, values map[interface{}]interface{}) error

----

#### HashHSETNX HSETNX key field value
> func (c *RedisClient) HashHSETNX(key, field string, value interface{}) error

----

#### HashHDEL HDEL key field [field ...] 删除哈希表
> func (c *RedisClient) HashHDEL(key string, fields []string) error

----

#### HashHEXISTS HEXISTS key field 查看哈希表
> func (c *RedisClient) HashHEXISTS(key, fields string) bool

----

#### HashHGET HGET key field 返回哈希表
> func (c *RedisClient) HashHGET(key, fields string) (string, error)

----

#### HashHINCRBY HINCRBY key field increment
> func (c *RedisClient) HashHINCRBY(key, field string, increment int64) (int64, error)

----

#### HashHINCRBYFLOAT HINCRBYFLOAT key field increment
> func (c *RedisClient) HashHINCRBYFLOAT(key, field string, increment float64) (res float64, err error)

----

#### HashHKEYS HKEYS key 返回哈希表
> func (c *RedisClient) HashHKEYS(key string) (res []string, err error)

----

#### HashHLEN HLEN key 返回哈希表
> func (c *RedisClient) HashHLEN(key string) (res int64, err error)

----

#### HashHMGET HMGET key field [field ...]
> func (c *RedisClient) HashHMGET(key string, fields []string) (res []string, err error)

----

#### HashHVALS HVALS key
> func (c *RedisClient) HashHVALS(key string) ([]string, error)

----

#### ListLRANGEALL LRANGE 获取List value
> func (c *RedisClient) ListLRANGEALL(key string) ([]interface{}, error)

----

#### ListLRANGE LRANGE key start stop
> func (c *RedisClient) ListLRANGE(key string, start, stop int64) ([]interface{}, error)

----

#### ListLPUSH LPUSH 新创建list 将一个或多个值 value 插入到列表 key 的表头
> func (c *RedisClient) ListLPUSH(key string, values []interface{}) error

----

#### ListRPUSH RPUSH key value [value ...]
> func (c *RedisClient) ListRPUSH(key string, values []interface{}) error

----

#### ListLINDEX LINDEX key index
> func (c *RedisClient) ListLINDEX(key string, index int64) (string, error)

----

#### ListLINSERT LINSERT key BEFORE|AFTER pivot value
> func (c *RedisClient) ListLINSERT(direction bool, key, pivot, value string) error

----

#### ListLLEN LLEN key
> func (c *RedisClient) ListLLEN(key string) (int64, error)

----

#### ListLPOP LPOP key
> func (c *RedisClient) ListLPOP(key string) (string, error)

----

#### ListLPUSHX LPUSHX key value
> func (c *RedisClient) ListLPUSHX(key string, value interface{}) error

----

#### ListLREM LREM key count value
> func (c *RedisClient) ListLREM(key string, count int64, value interface{}) error

----

#### ListLSET LSET key index value
> func (c *RedisClient) ListLSET(key string, index int64, value interface{}) error

----

#### ListLTRIM LTRIM key start stop
> func (c *RedisClient) ListLTRIM(key string, start, stop int64) error

----

#### ListRPOP RPOP key
> func (c *RedisClient) ListRPOP(key string) (string, error)

----

#### ListRPOPLPUSH RPOPLPUSH source destination
> func (c *RedisClient) ListRPOPLPUSH(key, destination string) (string, error)

----

#### ListRPUSHX RPUSHX key value
> func (c *RedisClient) ListRPUSHX(key string, value interface{}) error

----

#### SetSMEMBERS SMEMBERS key 返回集合 key 中的所有成员。
> func (c *RedisClient) SetSMEMBERS(key string) ([]interface{}, error)

----

#### SetSADD SADD 新创建Set  将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
> func (c *RedisClient) SetSADD(key string, values []interface{}) error

----

#### SetSCARD SCARD key
> func (c *RedisClient) SetSCARD(key string) error

----

#### SetSDIFF SDIFF key [key ...]
> func (c *RedisClient) SetSDIFF(keys []string) ([]interface{}, error)

----

#### SetSDIFFSTORE SDIFFSTORE destination key [key ...]
> func (c *RedisClient) SetSDIFFSTORE(key string, keys []string) ([]interface{}, error)

----

#### SetSINTER SINTER key [key ...]
> func (c *RedisClient) SetSINTER(keys []string) ([]interface{}, error)

----
#### SetSINTERSTORE SINTERSTORE destination key [key ...]
> func (c *RedisClient) SetSINTERSTORE(key string, keys []string) ([]interface{}, error)

----

#### SetSISMEMBER SISMEMBER key member
> func (c *RedisClient) SetSISMEMBER(key string, value interface{}) (resBool bool, err error)

----

#### SetSMOVE SMOVE source destination member
> func (c *RedisClient) SetSMOVE(key, destination string, member interface{}) (resBool bool, err error)

----

#### SetSPOP SPOP key
> func (c *RedisClient) SetSPOP(key string) (string, error)

----

#### SetSRANDMEMBER SRANDMEMBER key [count]
> func (c *RedisClient) SetSRANDMEMBER(key string, count int64) ([]interface{}, error)

----

#### SetSREM SREM key member [member ...]
> func (c *RedisClient) SetSREM(key string, member []interface{}) error

----

#### SetSUNION SUNION key [key ...]
> func (c *RedisClient) SetSUNION(keys []string) ([]interface{}, error)

----

#### SetSUNIONSTORE SUNIONSTORE destination key [key ...]
> func (c *RedisClient) SetSUNIONSTORE(key string, keys []string) ([]interface{}, error)

----

#### StringGet GET 获取String value
> func (c *RedisClient) StringGet(key string) (string, error)

----

#### StringSET SET 新建String
> func (c *RedisClient) StringSET(key string, value interface{}) error

----

#### StringSETEX SETEX 新建String 含有时间
> func (c *RedisClient) StringSETEX(key string, ttl int64, value interface{}) error

----

#### StringPSETEX PSETEX key milliseconds value
> func (c *RedisClient) StringPSETEX(key string, ttl int64, value interface{}) error

----

#### StringSETNX key value
> func (c *RedisClient) StringSETNX(key string, value interface{}) error

----
#### StringSETRANGE SETRANGE key offset value
> func (c *RedisClient) StringSETRANGE(key string, offset int64, value interface{}) error

----

#### StringAPPEND APPEND key value
> func (c *RedisClient) StringAPPEND(key string, value interface{}) error

----

#### StringDECR key
> func (c *RedisClient) StringDECR(key string) (int64, error)

----

#### StringDECRBY DECRBY key decrement
> func (c *RedisClient) StringDECRBY(key, decrement string) (int64, error)

----

#### StringGETRANGE GETRANGE key start end
> func (c *RedisClient) StringGETRANGE(key string, start, end int64) (string, error)

----

#### StringGETSET GETSET key value
> func (c *RedisClient) StringGETSET(key string, value interface{}) (string, error)

----
#### StringINCR INCR key
> func (c *RedisClient) StringINCR(key string) (int64, error)

----

#### StringINCRBY INCRBY key increment
> func (c *RedisClient) StringINCRBY(key, increment string) (int64, error)

----

#### StringINCRBYFLOAT INCRBYFLOAT key increment
> func (c *RedisClient) StringINCRBYFLOAT(key, increment float64) (float64, error)

----

#### StringMGET MGET key [key ...]
> func (c *RedisClient) StringMGET(key []interface{}) ([]string, error)

----

#### StringMSET MSET key value [key value ...]
> func (c *RedisClient) StringMSET(values []interface{}) error

---

#### StringMSETNX MSETNX key value [key value ...]
> func (c *RedisClient) StringMSETNX(values []interface{}) error

---


#### ZSetZRANGEALL ZRANGE 获取ZSet value 返回集合 有序集成员的列表。
> func (c *RedisClient) ZSetZRANGEALL(key string) ([]interface{}, error)

---


#### ZSetZRANGE ZRANGE key start stop [WITHSCORES]
> func (c *RedisClient) ZSetZRANGE(key string, start, stop int64) ([]interface{}, error)

---


#### ZSetZREVRANGE ZREVRANGE key start stop [WITHSCORES]
> func (c *RedisClient) ZSetZREVRANGE(key string, start, stop int64) ([]interface{}, error)

---


#### ZSetZADD ZADD 新创建ZSet 将一个或多个 member 元素及其 score 值加入到有序集 key 当中
> func (c *RedisClient) ZSetZADD(key string, values []interface{}) error

---


#### ZSetZCARD ZCARD key
> func (c *RedisClient) ZSetZCARD(key string) (int64, error)

---


#### ZSetZCOUNT ZCOUNT key min max
> func (c *RedisClient) ZSetZCOUNT(key string, min, max int64) (int64, error)

---


#### ZSetZINCRBY ZINCRBY key increment member
> func (c *RedisClient) ZSetZINCRBY(key, member string, increment int64) (string, error)

---


#### ZSetZRANGEBYSCORE ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
> func (c *RedisClient) ZSetZRANGEBYSCORE(key string, min, max, offset, count int64) ([]interface{}, error)

---


#### ZSetZRANGEBYSCOREALL ZRANGEBYSCORE
> func (c *RedisClient) ZSetZRANGEBYSCOREALL(key string) ([]interface{}, error)

---

#### ZSetZREVRANGEBYSCORE ZREVRANGEBYSCORE
> func (c *RedisClient) ZSetZREVRANGEBYSCORE(key string, min, max, offset, count int64) ([]interface{}, error)

---


#### ZSetZREVRANGEBYSCOREALL ZREVRANGEBYSCOREALL
> func (c *RedisClient) ZSetZREVRANGEBYSCOREALL(key string) ([]interface{}, error)

---


#### ZSetZRANK ZRANK key member
> func (c *RedisClient) ZSetZRANK(key string, member interface{}) (int64, error)

---


#### ZSetZREM ZREM key member [member ...]
> func (c *RedisClient) ZSetZREM(key string, member []interface{}) error

---


#### ZSetZREMRANGEBYRANK ZREMRANGEBYRANK key start stop
> func (c *RedisClient) ZSetZREMRANGEBYRANK(key string, start, stop int64) error

---

#### ZSetZREMRANGEBYSCORE ZREMRANGEBYSCORE key min max
> func (c *RedisClient) ZSetZREMRANGEBYSCORE(key string, min, max int64) error

---

#### ZSetZREVRANK ZREVRANK key member
> func (c *RedisClient) ZSetZREVRANK(key string, member interface{}) (int64, error)

---

#### ZSetZSCORE ZSCORE key member
> func (c *RedisClient) ZSetZSCORE(key string, member interface{}) (string, error)

---






