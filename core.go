package tinylog

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	OutFormat     = 0
	DefaultPrefix = time.Now().Format("2006/01/02 15:04:05")
	OutputPrefix  = time.Now().Format("2006-01-02 15-04-05")
)

type Logger struct {
	Prefix string // 日志前缀，用于标识日志的类型或来源
	mode   int    // 日志格式，如JSON，指定日志输出的结构
	count  int    // 日志缓冲区中允许存储的最大消息数量
	queue  []LogMsg
	mu     sync.Mutex // 互斥锁，确保并发写入日志时的线程安全
	wg     sync.WaitGroup
	Output *os.File // 日志文件输出流，用于将日志写入到文件
}

type LogMsg struct {
	time  string
	level int
	msg   string
}

var DefaultLogger = Logger{
	Prefix: DefaultPrefix,
	mode:   0,
	count:  0,
	queue:  make([]LogMsg, 0),
	mu:     sync.Mutex{},
	wg:     sync.WaitGroup{},
	Output: nil,
}

// New 方法用于创建一个新的Logger实例，并初始化其输出文件。
func New() *Logger {
	// 设置DefaultLogger的输出文件为新创建的文件
	DefaultLogger.Output = file()

	// 返回DefaultLogger的指针
	return &DefaultLogger
}

// Log 方法用于记录一条日志消息，支持不同的日志级别。
// level 参数指定日志的级别，msg 是要记录的消息内容。
func (l *Logger) Log(level int, msg string) {
	mid := ""
	// 根据日志级别设置中间前缀
	switch level {
	case 0:
		mid = " [Inf]"
	case 1:
		mid = " [Err]"
	default:
		mid = " [Unk]"
	}

	// 构造完整的日志消息
	logtext := l.Prefix + mid + msg + "\n"
	log.Print(logtext) // 将日志消息直接打印到控制台

	logmsg := &LogMsg{
		time:  DefaultPrefix,
		level: level,
		msg:   msg,
	}
	l.queue = append(l.queue, *logmsg)
	l.count++ // 更新缓冲区中的消息计数

	// 如果缓冲区中的日志消息数量达到或超过10条，调用Flush方法将它们写入文件
	if l.count >= 10 {
		l.wg.Add(1)
		go l.Flush()
	}
}

// Flush 方法用于将Logger缓冲区中的所有日志消息写入到日志文件中。
func (l *Logger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch l.mode {
	case 0:
		l.writeLog()
	case 1:
		l.writeJSON()
	default:
		log.Println("Flush时mode出错")
	}
	l.wg.Wait()
}
