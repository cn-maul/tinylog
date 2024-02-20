package tinylog

import (
	"log"
	"os"
	"strings"
	"time"
)

var DefaultPrefix = time.Now().Format("2006/01/02 15:04:05")

type Logger struct {
	Output *os.File
	Prefix string           //前缀
	Format string           //支持结构化日志（如JSON格式）
	Buf    *strings.Builder //暂时存储日志队列
}

var DefaultLogger = Logger{
	Output: nil,
	Prefix: DefaultPrefix,
	Format: "",
	Buf:    new(strings.Builder),
}

func New() *Logger {
	output, err := os.Open("./log/" + DefaultPrefix + ".log")
	if err != nil {
		log.Println(err)
	}

	DefaultLogger.Output = output
	return &DefaultLogger
}
