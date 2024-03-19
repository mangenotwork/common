package ginHelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type CommonResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code  int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`  // 错误代码
	Data  []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`   // 成功响应数据
	Msg   string `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`     // 消息
	Stack string `protobuf:"bytes,4,opt,name=stack,proto3" json:"stack,omitempty"` // 服务器内部错误时的调用栈（用于开发环境的调试使用）
}

type GrpcErrorCode int32

const (
	ErrorCode_NO               GrpcErrorCode = 0
	ErrorCode_UnKnow           GrpcErrorCode = 1 //未知错误
	ErrorCode_InternalServer   GrpcErrorCode = 2 //内部服务错误
	ErrorCode_RequestData      GrpcErrorCode = 3 //请求数据错误
	ErrorCode_InterfaceInvalid GrpcErrorCode = 4 //接口无效
	ErrorCode_ServerStop       GrpcErrorCode = 5 //服务已暂停
	ErrorCode_ErrDBNoRecord    GrpcErrorCode = 6 //没有记录
	ErrorCode_Database         GrpcErrorCode = 7 //数据错误
)

func RespWriteGrpcOK(ctx *gin.Context, data interface{}) {
	var temp []byte
	if protoData, ok := data.(protoreflect.ProtoMessage); ok {
		temp, _ = proto.Marshal(protoData)
	}
	result := &CommonResp{
		Code: int32(ErrorCode_NO),
		Msg:  "Ok",
		Data: temp,
	}
	ctx.ProtoBuf(http.StatusOK, result)
}

func RespWriteGrpcFail(ctx *gin.Context, code GrpcErrorCode, err error) {
	result := &CommonResp{
		Code: int32(code),
		Msg:  err.Error(),
	}
	ctx.ProtoBuf(http.StatusOK, result)
}

func ReqGrpcParameter(ctx *gin.Context, obj any) error {
	return ctx.MustBindWith(obj, binding.ProtoBuf)
}
