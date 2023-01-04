package log

import (
	"testing"
)

// go test -run Test_log_out
// go test -v -bench=. logger_test.go logger.go
// go test -v -bench=. -benchtime=5s logger_test.go logger.go
func Test_log_out(t *testing.T) {
	Print("这是日志打印")
	Info("这是info打印")
	InfoF("这是infoF打印 %s", "infoF")
	InfoTimes(0, "这是InfoTimes打印")
	Debug("这是Debug打印")
	DebugF("这是DebugF打印 %s", "DebugF")
	DebugTimes(0, "这是DebugTimes打印")
	Warn("这是Warn打印")
	WarnF("这是WarnF打印")
	WarnTimes(0, "这是WarnTimes打印")
	Error("这是Error打印")
	ErrorF("这是ErrorF打印")
	ErrorTimes(0, "这是ErrorTimes打印")
	//Panic("这是Panic打印")
}

// go test -run Test_log_file
func Test_log_file(t *testing.T) {
	SetLogFile("")
	Print("这是日志打印")
	Info("这是info打印")
	InfoF("这是infoF打印 %s", "infoF")
	InfoTimes(0, "这是InfoTimes打印")
	Debug("这是Debug打印")
}

// go test -run Test_log_out_service
func Test_log_out_service(t *testing.T) {
	SetOutService("127.0.0.1", 8222)
	Print("这是日志打印")
	Info("这是info打印")
	InfoF("这是infoF打印 %s", "infoF")
	InfoTimes(0, "这是InfoTimes打印")
	Debug("这是Debug打印")
}
