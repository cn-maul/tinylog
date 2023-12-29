package tlog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02_15-04-05"
)

// 日志的结构体
type Log struct {
	Start   string          //记录日志的开始时间,作为文件名
	Content strings.Builder //日志内容
}

// 单条日志信息的结构体
type Logger struct {
	Level   string //日志等级
	Content string //日志内容
}

// New 创建一个新的日志对象
func New() *Log {
	return &Log{
		Start:   time.Now().Format(timeFormat),
		Content: strings.Builder{},
	}
}

// Write 将日志写入本地文件
func (l *Log) Write(path string) {
	// 检查路径是否存在
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(err)
		}
		return
	}

	// 创建一个文件
	filename := l.Start + ".log"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 将数据写入文件
	_, err = io.WriteString(f, l.Content.String())
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Add 将一条日志信息添加到日志对象中
func (l *Log) Add(g *Logger) error {
	log := "[" + g.Level + "]" + time.Now().Format(timeFormat) + ":" + g.Content + "\n"
	l.Content.Grow(len(log))
	_, err := l.Content.WriteString(log)
	return err
}

// Info 记录一条信息日志
func (l *Log) Info(content string) {
	err := l.Add(&Logger{
		Level:   "Info ",
		Content: content,
	})
	if err != nil {
		fmt.Println("Error Writing Builder")
	}
}

// Error 记录一条错误日志
func (l *Log) Error(content string) {
	err := l.Add(&Logger{
		Level:   "Error",
		Content: content,
	})
	if err != nil {
		fmt.Println("Error Writing Builder")
	}
}

// Fatal 记录一条致命日志
func (l *Log) Fatal(content string) {
	err := l.Add(&Logger{
		Level:   "Fatal",
		Content: content,
	})
	if err != nil {
		fmt.Println("Error Writing Builder")
	}
}
