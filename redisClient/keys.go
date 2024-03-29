package redisClient

import (
	"fmt"
	"strings"

	"github.com/mangenotwork/common/log"

	"github.com/garyburd/redigo/redis"
)

// GetALLKeys 获取所有的key
func (c *RedisClient) GetALLKeys(match string) (ksyList map[string]int) {

	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil
	}

	//初始化拆分值
	matchSplit := match
	//match :匹配值，没有则匹配所有 *
	if match == "" {
		match = "*"
	} else {
		match = fmt.Sprintf("*%s*", match)
	}
	//cursor :初始游标为0
	cursor := "0"
	ksyList = make(map[string]int)
	ksyList, cursor = c.addGetKey(conn, ksyList, cursor, match, matchSplit)
	//当游标等于0的时候停止获取key
	//线性获取，一直循环获取key,直到游标为0
	if cursor != "0" {
		for {
			ksyList, cursor = c.addGetKey(conn, ksyList, cursor, match, matchSplit)
			if cursor == "0" {
				break
			}
		}
	}
	return
}

// addGetKey 内部方法
// 针对分组的key进行分组合并处理
func (c *RedisClient) addGetKey(conn redis.Conn, ksyList map[string]int, cursor, match, matchSplit string) (map[string]int, string) {

	countNumber := "10000"
	res, err := redis.Values(conn.Do("scan", cursor, "MATCH", match, "COUNT", countNumber))
	log.InfoTimes(3, "[Redis Log] execute :", "scan ", cursor, " MATCH ", match, " COUNT ", countNumber)
	if err != nil {
		log.Error("GET error", err.Error())
	}
	//获取	match 含有多少:
	cfNumber := strings.Count(match, ":")
	//获取新的游标
	newCursor := string(res[0].([]byte))
	allKey := res[1]
	allKeyData := allKey.([]interface{})
	for _, v := range allKeyData {
		keyData := string(v.([]byte))
		//没有:的key 则不集合
		if strings.Count(keyData, ":") == cfNumber || keyData == match {
			ksyList[keyData] = 0
			continue
		}
		//有:需要集合
		keyDataNew, _ := fenGeYinHaoOne(keyData, matchSplit)
		ksyList[keyDataNew] = ksyList[keyDataNew] + 1
	}
	return ksyList, newCursor
}

// fenGeYinHaoOne 对查询出来的key进行拆分，集合，分组处理
func fenGeYinHaoOne(str string, matchSplit string) (string, int) {
	likeKey := ""
	if matchSplit != "" {
		likeKey = fmt.Sprintf("%s", matchSplit)
	}
	str = strings.Replace(str, likeKey, "", 1)
	fg := strings.Split(str, ":")
	if len(fg) > 0 {
		return fmt.Sprintf("%s%s", likeKey, fg[0]), len(fg)
	}
	return "", len(fg)
}

func (c *RedisClient) SearchKeys(match string) (ksyList map[string]int) {

	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil
	}

	ksyList = make(map[string]int)
	if match == "" {
		return
	} else {
		match = fmt.Sprintf("*%s*", match)
	}
	cursor := "0"
	ksyList = make(map[string]int)
	ksyList, cursor = c.addSearchKey(conn, ksyList, cursor, match)
	//当游标等于0的时候停止获取key
	//线性获取，一直循环获取key,直到游标为0
	if cursor != "0" {
		for {
			ksyList, cursor = c.addSearchKey(conn, ksyList, cursor, match)
			if cursor == "0" {
				break
			}
		}
	}
	return
}

// addGetKey 内部方法获取key
func (c *RedisClient) addSearchKey(conn redis.Conn, ksyList map[string]int, cursor, match string) (map[string]int, string) {
	countNumber := "10000"
	res, err := redis.Values(conn.Do("scan", cursor, "MATCH", match, "COUNT", countNumber))
	log.InfoTimes(3, "[Redis Log] execute :", "scan ", cursor, " MATCH ", match, " COUNT ", countNumber)
	if err != nil {
		log.Error("GET error", err.Error())
	}
	//获取新的游标
	newCursor := string(res[0].([]byte))
	allKey := res[1]
	allKeyData := allKey.([]interface{})
	for _, v := range allKeyData {
		keyData := string(v.([]byte))
		ksyList[keyData] = 0
	}
	return ksyList, newCursor
}

// GetKeyType 获取key的类型
func (c *RedisClient) GetKeyType(key string) string {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return ""
	}
	log.InfoFTimes(3, "[Redis Log] execute : TYPE %s", key)
	res, err := redis.String(conn.Do("TYPE", key))
	if err != nil {
		log.ErrorTimes(3, "GET error", err.Error())
	}
	return res
}

// GetKeyTTL 获取key的过期时间
func (c *RedisClient) GetKeyTTL(key string) int64 {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return 0
	}
	log.InfoFTimes(3, "[Redis Log] execute : TTL %s", key)
	res, err := redis.Int64(conn.Do("TTL", key))
	if err != nil {
		log.ErrorTimes(3, "GET error", err.Error())
	}
	return res
}

// EXISTSKey 检查给定 key 是否存在。
func (c *RedisClient) EXISTSKey(key string) bool {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return false
	}
	log.InfoFTimes(3, "[Redis Log] execute : DUMP %s", key)
	data, err := redis.String(conn.Do("DUMP", key))
	if err != nil || data == "0" {
		log.ErrorTimes(3, "GET error", err.Error())
		return false
	}
	return true
}

// RenameKey 修改key名称
func (c *RedisClient) RenameKey(name, newName string) bool {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return false
	}
	arg := redis.Args{}.Add(name).Add(newName)
	log.InfoFTimes(3, "[Redis Log] execute : RENAME %s %v", name, newName)
	_, err = conn.Do("RENAME", arg...)
	if err != nil {
		log.ErrorTimes(3, "GET error", err.Error())
		return false
	}
	return true
}

// UpdateKeyTTL 更新key ttl
func (c *RedisClient) UpdateKeyTTL(key string, ttl int64) bool {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return false
	}
	arg := redis.Args{}.Add(key).Add(ttl)
	log.InfoFTimes(3, "[Redis Log] execute : EXPIRE %s %v", key, ttl)
	_, err = conn.Do("EXPIRE", arg...)
	if err != nil {
		log.ErrorTimes(3, err.Error())
		return false
	}
	return true
}

// EXPIREATKey 指定key多久过期 接收的是unix时间戳
func (c *RedisClient) EXPIREATKey(key string, date int64) bool {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return false
	}
	arg := redis.Args{}.Add(key).Add(date)
	log.InfoFTimes(3, "[Redis Log] execute : EXPIREAT %s %v", key, date)
	_, err = conn.Do("EXPIREAT", arg...)
	if err != nil {
		log.ErrorTimes(3, err.Error())
		return false
	}
	return true
}

// DELKey 删除key
func (c *RedisClient) DELKey(key string) bool {
	conn, err := GetConn(c.Name)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return false
	}
	log.InfoFTimes(3, "[Redis Log] execute : DEL %s", key)
	_, err = conn.Do("DEL", key)
	if err != nil {
		log.ErrorTimes(3, err.Error())
		return false
	}
	return true
}
