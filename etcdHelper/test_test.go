package etcdHelper

import (
	"github.com/mangenotwork/common/log"
	"testing"
	"time"
)

var etcdAddr = []string{"127.0.0.1:2379"}

func Test_Base(t *testing.T) {
	c, err := Conn(etcdAddr, 3*time.Second)
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		_ = c.Close()
	}()
	k := "test-k"
	//v := "test-v"
	//err = Put(c, k, v)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	resp, err := Get(c, k)
	if err != nil {
		log.Error(err)
		return
	}
	for a, b := range resp.Kvs {
		log.Debug(a, "-> ", string(b.Key), ":", string(b.Value))
	}
}

// 注册服务
func Test_Register(t *testing.T) {
	register := NewServiceRegister(etcdAddr)
	_, err := register.Register(Server{
		Name:     "test-1216",
		Addr:     "127.0.0.1:12161",
		Version:  "1",
		Metadata: nil,
	}, 10)
	if err != nil {
		log.Error(err)
	}
}

// 发现服务
func Test_Discovery(t *testing.T) {
}
