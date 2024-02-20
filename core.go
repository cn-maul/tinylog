package tinylog

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	InfoLevel  = "inf"
	ErrorLevel = "err"
)

var DefaultPrefix = time.Now().Format("2006/01/02 15:04:05")
var OutputPrefix = time.Now().Format("2006-01-02 15-04-05")

type Logger struct {
	Prefix string           // 日志前缀，用于标识日志的类型或来源
	Format string           // 日志格式，如JSON，用于指定日志输出的结构
	BufNum int              // 日志缓冲区中允许存储的最大消息数量
	Buf    *strings.Builder // 日志缓冲区，用于暂存日志消息，以优化性能
	mu     sync.Mutex       // 互斥锁，确保并发写入日志时的线程安全
	Output *os.File         // 日志文件输出流，用于将日志写入到文件
}

var DefaultLogger = Logger{
	Prefix: DefaultPrefix,
	Format: "",
	BufNum: 0,
	Buf:    new(strings.Builder),
	mu:     sync.Mutex{},
	Output: nil,
}

// New 方法用于创建一个新的Logger实例，并初始化其输出文件。
func New() *Logger {
	checkOutput() // 检查并创建日志输出目录

	// 创建日志文件，文件名由当前时间的日期格式和".log"后缀组成
	output, err := os.Create("./log/" + OutputPrefix + ".log")
	if err != nil { // 如果创建文件失败，记录错误信息
		log.Println(err)
	}

	// 设置DefaultLogger的输出文件为新创建的文件
	DefaultLogger.Output = output

	// 返回DefaultLogger的指针
	return &DefaultLogger
}

// Log 方法用于记录一条日志消息，支持不同的日志级别。
// level 参数指定日志的级别，msg 是要记录的消息内容。
func (l *Logger) Log(level, msg string) {
	// 根据日志级别设置中间前缀
	midPrefix := ""
	switch level {
	case InfoLevel:
		midPrefix = " [Inf]"
	case ErrorLevel:
		midPrefix = " [Err]"
	default:
		midPrefix = " [Unk]"
	}

	// 如果缓冲区中的日志消息数量达到或超过10条，调用Flush方法将它们写入文件
	if l.BufNum >= 10 {
		l.Flush()
	}

	// 构造完整的日志消息
	logMsg := l.Prefix + midPrefix + msg + "\n"
	log.Print(logMsg) // 将日志消息直接打印到控制台

	// 将日志消息追加到缓冲区
	l.Buf.Grow(len(logMsg))   // 确保缓冲区有足够的空间
	l.Buf.WriteString(logMsg) // 写入日志消息
	l.BufNum++                // 更新缓冲区中的消息计数
}

// Flush 方法用于将Logger缓冲区中的所有日志消息写入到日志文件中。
func (l *Logger) Flush() {
	l.mu.Lock()         // 获取互斥锁，确保并发安全
	defer l.mu.Unlock() // 释放互斥锁

	// 尝试将缓冲区的内容写入日志文件
	_, err := l.Output.WriteString(l.Buf.String())
	if err != nil { // 如果写入失败，记录错误信息
		log.Println(err)
	}

	// 重置缓冲区，准备接收新的日志消息
	l.reset()
}

// reset 方法用于重置Logger的缓冲区，清除所有已存储的日志消息。
func (l *Logger) reset() {
	l.BufNum = 0  // 重置缓冲区计数器
	l.Buf.Reset() // 清空缓冲区内容
}
