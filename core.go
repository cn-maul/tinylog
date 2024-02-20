package tinylog

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var DefaultPrefix = time.Now().Format("2006/01/02 15:04:05")
var OutputPrefix = time.Now().Format("2006-01-02 15-04-05")

type Logger struct {
	Output *os.File
	Prefix string           //前缀
	Format string           //支持结构化日志（如JSON格式）
	Buf    *strings.Builder //暂时存储日志队列
	BufNum int              //对列允许存储的最大消息数量计数器
}

var DefaultLogger = Logger{
	Output: nil,
	Prefix: DefaultPrefix,
	Format: "",
	Buf:    new(strings.Builder),
	BufNum: 0,
}

func New() *Logger {
	checkOutput()
	output, err := os.Create("./log/" + OutputPrefix + ".log")
	if err != nil {
		log.Println(err)
	}

	DefaultLogger.Output = output
	return &DefaultLogger
}

func (l *Logger) reset() {
	l.BufNum = 0
	l.Buf.Reset()
}

func (l *Logger) Println(logLevel, msg string) {
	var midPrefix string
	switch logLevel {
	case "log":
		midPrefix = " [Msg]"
	case "err":
		midPrefix = " [Err]"
	default:
		midPrefix = " [Nil]"
	}

	if l.BufNum == 10 {
		//写入文件的函数
		l.Write()
		l.reset()
	}

	logMsg := l.Prefix + midPrefix + msg + "\n"
	fmt.Print(logMsg)
	l.Buf.Grow(len(logMsg))
	l.Buf.WriteString(logMsg)
}
