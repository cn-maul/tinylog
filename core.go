package tlog

import (
	"os"
	"os/signal"
	"syscall"
)

// Log 日志结构体
type Log struct {
	Time    string // 时间
	Level   int    // 日志级别
	Message string // 日志消息
}

// init 初始化函数
func init() {
	// 创建日志目录
	os.MkdirAll(OutputDir, 0777)

	// 启动后台服务
	go backService(LogChan)

	// 启动fatalWrite goroutine
	go fatalWrite()
}

// fatalWrite 捕获SIGINT和SIGTERM信号，写入日志并退出程序
func fatalWrite() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range c {
			write()    // 写入日志
			os.Exit(0) // 退出程序
		}
	}()
}
