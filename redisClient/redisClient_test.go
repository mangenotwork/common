package redisClient

import (
	"github.com/mangenotwork/common/log"
	"testing"
)

// go test -test.run Test_client_conn -v
func Test_client_conn(t *testing.T) {
	c1, err1 := Conn("127.0.0.1:6379", "123", "7")
	log.Print(c1, err1)
	err2 := RedisConn("127.0.0.1:6379", "123", "0")
	log.Print(Client, err2)

}

// go test -test.run Test_hash -v
func Test_hash(t *testing.T) {
	err2 := RedisConn("127.0.0.1:6379", "123", "0")
	log.Print(Client, err2)
	t1, e1 := Client.HashHSET("a1", "1", "111")
	log.Print(t1, e1)
	t2, e2 := Client.HashHGETALL("a1")
	log.Print(t2, e2)
	log.Print(Client.HashHEXISTS("a1", "1"))
	log.Print(Client.HashHEXISTS("a1", "2"))
	log.Print(Client.HashHEXISTS("a2", "2"))
	log.Print(Client.HashHMSET("a2", map[interface{}]interface{}{1: "a", 2: "b"}))
	log.Print(Client.HashHGETALL("a2"))
}
