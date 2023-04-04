package log

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var std = newStd()

type logger struct {
	appName         string
	terminal        bool
	outFile         bool
	outFileWriter   *os.File
	outService      bool // 日志输出到服务
	outServiceIp    string
	outServicePort  int
	outServiceConn  *net.UDPConn
	outServiceLevel []int
}

func newStd() *logger {
	return &logger{
		terminal:        true,
		outFile:         false,
		outService:      false,
		outServiceLevel: []int{3, 4, 5},
	}
}

// SetLogFile 设置日志文件名称, 文件名称可含路径(绝对或相对)
func SetLogFile(name string) {
	std.outFile = true
	std.appName = name
	lastName := strings.Split(name, "/")
	if len(lastName) > 0 && len(lastName[len(lastName)-1]) > 0 {
		name = name + "_"
	}
	std.outFileWriter, _ = os.OpenFile(name+time.Now().Format("20060102")+".log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// SetAppName 设置项目名称
func SetAppName(name string) {
	std.appName = name
}

func SetOutService(ip string, port int) {
	var err error
	std.outService = true
	std.outServiceIp = ip
	std.outServicePort = port
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: net.ParseIP(std.outServiceIp), Port: std.outServicePort}
	std.outServiceConn, err = net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		Error(err)
	}
}

func SetOutServiceWarn2Panic() {
	std.outServiceLevel = []int{3, 4, 5}
}

func SetOutServiceInfo2Panic() {
	std.outServiceLevel = []int{1, 2, 3, 4, 5}
}

func Close() {
	std.terminal = false
	std.outFile = false
	std.outService = false
}

func DisableTerminal() {
	std.terminal = false
}

type Level int

var LevelMap = map[Level]string{
	0: "[Print] ",
	1: "[INFO]  ",
	2: "[DEBUG] ",
	3: "[WARN]  ",
	4: "[ERROR] ",
	5: "[PANIC] ",
	6: "[Http]",
}

func (l *logger) Log(level Level, args string, times int) {
	var buffer bytes.Buffer
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05.000 "))
	buffer.WriteString(LevelMap[level])
	_, file, line, _ := runtime.Caller(times)
	fileList := strings.Split(file, "/")
	// 最多显示两级路径
	if len(fileList) > 3 {
		fileList = fileList[len(fileList)-3 : len(fileList)]
	}

	if times != -1 {
		buffer.WriteString(strings.Join(fileList, "/"))
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(line))
	}

	buffer.WriteString(" \t| ")
	buffer.WriteString(args)
	buffer.WriteString("\n")
	out := buffer.Bytes()
	if l.terminal {
		_, _ = buffer.WriteTo(os.Stdout)
	}
	// 输出到文件或远程日志服务
	if l.outFile {
		_, _ = l.outFileWriter.Write(out)
	}
	if l.outService {
		for _, v := range l.outServiceLevel {
			if Level(v) == level {
				out = append([]byte("1"+l.appName+"|"), out...)
				_, _ = l.outServiceConn.Write(out)
			}
		}
	}
}

func Print(args ...interface{}) {
	std.Log(0, fmt.Sprint(args...), 2)
}

func PrintF(format string, args ...interface{}) {
	std.Log(0, fmt.Sprintf(format, args...), 2)
}

func Info(args ...interface{}) {
	std.Log(1, fmt.Sprint(args...), 2)
}

func InfoF(format string, args ...interface{}) {
	std.Log(1, fmt.Sprintf(format, args...), 2)
}

// InfoTimes
// times 意思是打印第几级函数调用
func InfoTimes(times int, args ...interface{}) {
	std.Log(1, fmt.Sprint(args...), times)
}

// InfoFTimes
// times 意思是打印第几级函数调用
func InfoFTimes(times int, format string, args ...interface{}) {
	std.Log(1, fmt.Sprintf(format, args...), times)
}

func Debug(args ...interface{}) {
	std.Log(2, fmt.Sprint(args...), 2)
}

func DebugF(format string, args ...interface{}) {
	std.Log(2, fmt.Sprintf(format, args...), 2)
}

func DebugTimes(times int, args ...interface{}) {
	std.Log(2, fmt.Sprint(args...), times)
}

func DebugFTimes(times int, format string, args ...interface{}) {
	std.Log(2, fmt.Sprintf(format, args...), times)
}

func Warn(args ...interface{}) {
	std.Log(3, fmt.Sprint(args...), 2)
}

func WarnF(format string, args ...interface{}) {
	std.Log(3, fmt.Sprintf(format, args...), 2)
}

func WarnTimes(times int, args ...interface{}) {
	std.Log(3, fmt.Sprint(args...), times)
}

func WarnFTimes(times int, format string, args ...interface{}) {
	std.Log(3, fmt.Sprintf(format, args...), times)
}

func Error(args ...interface{}) {
	std.Log(4, fmt.Sprint(args...), 2)
}

func ErrorF(format string, args ...interface{}) {
	std.Log(4, fmt.Sprintf(format, args...), 2)
}

func ErrorTimes(times int, args ...interface{}) {
	std.Log(4, fmt.Sprint(args...), times)
}

func ErrorFTimes(times int, format string, args ...interface{}) {
	std.Log(4, fmt.Sprintf(format, args...), times)
}

func Panic(args ...interface{}) {
	std.Log(5, fmt.Sprint(args...), 2)
	panic(args)
}

func HttpInfo(args ...interface{}) {
	std.Log(6, fmt.Sprint(args...), -1)
}
