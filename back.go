package tlog

import (
	"log"
	"time"
)

// Cubage 日志缓存上限
var Cubage = 100

// LogChan 日志通道
var LogChan = make(chan *Log, 1024)

// List 日志列表
var List []Log

// 日志的后台服务
func backService(logChan chan *Log) {
	// 从channel中获取对象
	for logger := range logChan {
		// 接受到退出信号，退出后台服务
		if (*logger).Level == -1 {
			log.Println("后端退出")
			return
		}

		List = append(List, *logger)
		if len(List) >= Cubage {
			write() // 写入日志
		}
	}
}

// Println 打印日志
func Println(msg string) {
	LogChan <- &Log{
		Time:    time.Now().Format("2006/01/02 15:04:05"),
		Level:   0,
		Message: msg,
	}
	log.Printf("[%s] %s\n", "Info", msg)
}

// Errorln 打印错误日志
func Errorln(msg string) {
	LogChan <- &Log{
		Time:    time.Now().Format("2006/01/02 15:04:05"),
		Level:   1,
		Message: msg,
	}
	log.Printf("[%s] %s\n", "Error", msg)
}
