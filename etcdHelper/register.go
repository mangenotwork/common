package etcdHelper

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mangenotwork/common/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

var DialTimeout = 3

// Register 注册服务信息
type Register struct {
	EtcdAddr    []string
	DialTimeout int
	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
	cli         *clientv3.Client
	srvInfo     Server
	srvTTL      int64
}

// NewRegister 新建注册服务
func NewRegister(cli *clientv3.Client) *Register {
	return &Register{
		cli:         cli,
		DialTimeout: DialTimeout,
	}
}

// NewServiceRegister 新建注册服务
func NewServiceRegister(etcdAddr []string) *Register {
	return &Register{
		EtcdAddr:    etcdAddr,
		DialTimeout: DialTimeout,
	}
}

// Register 注册服务
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error

	if nil == r.cli {
		if strings.Split(srvInfo.Addr, ":")[0] == "" {
			return nil, errors.New("无效IP")
		}

		if r.cli, err = clientv3.New(clientv3.Config{
			Endpoints:   r.EtcdAddr,
			DialTimeout: time.Duration(r.DialTimeout) * time.Second,
		}); err != nil {
			return nil, err
		}
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})

	go r.keepAlive()

	return r.closeCh, nil
}

// Stop 停止注册服务
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// GetServerInfo 获取现有的服务信息
func (r *Register) GetServerInfo() (Server, error) {
	resp, err := r.cli.Get(context.Background(), BuildRegPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}
	info := Server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}
	return info, nil
}

// register 注册节点
func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// 	设置租约时间
	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}
	r.leasesID = leaseResp.ID
	// 设置续租 定期发送请求
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	// 注册服务并绑定租约
	_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

// unregister 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegPath(r.srvInfo))
	return err
}

// keepAlive 设置续租
func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				log.Error("删除失败", zap.Error(err))
			}
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				log.Error("撤销失败", zap.Error(err))
			}
			return
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					log.Error("注册失败", zap.Error(err))
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					log.Error("注册失败", zap.Error(err))
				}
			}
		}
	}
}
