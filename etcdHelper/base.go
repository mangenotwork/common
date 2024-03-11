package etcdHelper

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

/*
	包含etcd基础操作
	连接，增删改查
*/

func Conn(endpoints []string, timeout time.Duration) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
}

func Put(c *clientv3.Client, k, v string) error {
	kvc := clientv3.NewKV(c)
	_, err := kvc.Put(context.Background(), k, v)
	return err
}

func Get(c *clientv3.Client, k string) (*clientv3.GetResponse, error) {
	kvc := clientv3.NewKV(c)
	return kvc.Get(context.Background(), k)
}
