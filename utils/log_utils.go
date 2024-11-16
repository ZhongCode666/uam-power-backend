package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	mutex     sync.Mutex
	logger    *log.Logger
	logDir    = "./logs" // 日志文件存放的目录
	logPrefix = "log-"   // 日志文件名前缀
)

// Log 结构体
type Log struct{}

// NewLog 创建一个新的 Log 实例
func NewLog() *Log {
	return &Log{}
}

// Init 初始化日志系统，设置日志文件
func (l *Log) Init() error {
	// 获取当前时间，并格式化为日志文件名
	//now := time.Now()
	//logFileName := fmt.Sprintf("%s/%s%s.log", logDir, logPrefix, now.Format("20060102150405"))

	// 确保日志目录存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	// 打开日志文件，文件不存在则创建
	//file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	//if err != nil {
	//	return fmt.Errorf("failed to open log file: %v", err)
	//}

	// 使用 log.New 创建一个 Logger
	//logger = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	return nil
}

// Write 写入日志
func (l *Log) Write(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	// 输出日志，带时间戳
	logger.Println(message)
}

func (l *Log) Info(message string) {
	l.Write(fmt.Sprintf("%s☐", message))
}

func (l *Log) Error(message string) {
	l.Write(fmt.Sprintf("%s☒", message))
}

func (l *Log) Success(message string) {
	l.Write(fmt.Sprintf("%s☑", message))
}

// Close 关闭日志文件
func (l *Log) Close() {
	if logger != nil && logger.Writer() != nil {
		if file, ok := logger.Writer().(*os.File); ok {
			_ = file.Close()
		}
	}
}

var LogInfo Log

func InitLog() {
	err := LogInfo.Init()
	if err != nil {
		return
	}
}

func MsgError(msg string) {
	LogInfo.Error(msg)
}

func MsgInfo(msg string) {
	LogInfo.Info(msg)
}

func MsgSuccess(msg string) {
	LogInfo.Success(msg)
}
