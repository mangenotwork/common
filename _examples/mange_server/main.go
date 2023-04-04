package main

import (
	mange "github.com/mangenotwork/common/httpMiddleware"
	"net/http"
)

func main() {
	mux := mange.NewEngine()
	mux.Router("/", Index)
	mux.Run("18888")
}

func Index(w http.ResponseWriter, r *http.Request) {
	mange.OutSucceedBodyJsonP(w, "ok")
}

// TODO : 1. 路由没法添加中间件
// TODO : 2. Handler写法太复杂
// TODO : 3. 日志输出的文件行号不对
