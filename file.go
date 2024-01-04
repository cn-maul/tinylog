package tlog

import (
	"fmt"
	"os"
	"time"
)

// LogLevel 定义日志级别和对应的字符串表示
var LogLevel = map[int]string{
	0: "Info",
	1: "Error",
}

// Logger 结构体表示日志器
type Logger struct {
	file *os.File
}

// NewLogger 创建一个新的Logger实例，并打开指定文件
func NewLogger() (*Logger, error) {
	timestamp := time.Now().Format("2006-01-02 15-04-05")
	file, err := os.Create(timestamp + ".log")
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

// Close 方法用于关闭日志文件
func (l *Logger) Close() {
	l.file.Close()
}

func (l *Logger) Info(format string, a ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := timestamp + " " + LogLevel[0] + " " + fmt.Sprintf(format, a) + "\n"
	fmt.Println(logMessage)
	l.file.WriteString(logMessage)
}

func (l *Logger) Error(err string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintln(timestamp, LogLevel[1], "["+err+"]")
	fmt.Print(logMessage)
	l.file.WriteString(logMessage)
}
