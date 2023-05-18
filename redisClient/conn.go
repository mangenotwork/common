package redisClient

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func Conn(host, password, db string) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", host)
	if nil != err {
		return nil, err
	}
	if _, authErr := conn.Do("AUTH", password); authErr != nil {
		return nil, fmt.Errorf("redis auth password error: %s", authErr)
	}
	_, err = conn.Do("select", db)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func RedisPoolConn(host, password, db string, maxIdle, maxActive int) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			if conn == nil {
				return nil, errors.New("连接redis错误")
			}
			if len(password) != 0 {
				if _, err = conn.Do("AUTH", password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			if _, err = conn.Do("SELECT", db); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
	return pool
}

func getSSHClient(user string, pass string, addr string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshConn, err := net.Dial("tcp", addr)
	if nil != err {
		return nil, err
	}
	clientConn, chans, reqs, err := ssh.NewClientConn(sshConn, addr, config)
	if nil != err {
		sshConn.Close()
		return nil, err
	}
	client := ssh.NewClient(clientConn, chans, reqs)
	return client, nil
}

func RedisSshConn(user, pass, addr, host, password, db string) (redis.Conn, error) {
	var err error
	sshClient, err := getSSHClient(user, pass, addr)
	if nil != err {
		return nil, err
	}
	conn, err := sshClient.Dial("tcp", host)
	if nil != err {
		return nil, err
	}
	newConn := redis.NewConn(conn, -1, -1)
	if _, authErr := newConn.Do("AUTH", password); authErr != nil {
		return nil, fmt.Errorf("redis auth password error: %s", authErr)
	}
	if db == "" {
		db = "0"
	}
	_, err = newConn.Do("select", db)
	if nil != err {
		return nil, err
	}
	return newConn, nil
}
