package grpcHelper

import (
	"net"

	"google.golang.org/grpc"
)

// SimpleServer 启动一个简单的grpc服务
func SimpleServer(addr string, g *grpc.Server) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic("监听错误:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("启动错误:" + err.Error())
	}
}
