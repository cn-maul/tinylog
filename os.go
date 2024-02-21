package tinylog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	Level = map[int]string{
		0: "[Inf]",
		1: "[Err]",
	}
	Mode = map[int]string{
		0: ".log",
		1: ".json",
	}
)

// Write 方法用于将Logger缓冲区中的日志消息写入到文件，并刷新缓冲区
func (l *Logger) write(buf string) {
	// 创建一个带有缓冲的写入器，提高写入效率
	writer := bufio.NewWriter(l.Output)

	// 尝试将缓冲区的内容写入文件
	_, err := writer.WriteString(buf)
	if err != nil {
		log.Println("Bufio写时出错:", err)
		return
	}

	// 刷新写入器，确保所有数据都被写入文件
	err = writer.Flush()
	if err != nil {
		log.Println("Flush时出错:", err)
		return
	}
}

// checkOutput 方法用于检查日志输出目录是否存在，并在需要时创建它。
func checkOutput() {
	logDir := "./log" // 日志输出目录的路径

	// 检查日志目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 如果目录不存在，创建它
		log.Println("日志文件夹不存在，正在创建log目录")
		// 创建目录并设置适当的权限
		if err = os.MkdirAll(logDir, 0o755); err != nil {
			log.Println("创建log目录失败:", err)
		}
	}
}

func (l *Logger) writeLog() {
	var buf strings.Builder

	for _, v := range l.queue {
		content := l.Prefix + Level[v.level] + v.msg + "\n"
		buf.Grow(len(content))
		buf.WriteString(content)
	}
	l.write(buf.String())

	// 重置缓冲区，准备接收新的日志消息
	l.count = 0
	l.queue = l.queue[:0]
	defer l.wg.Done()
}

func (l *Logger) writeJSON() {
	// 将切片转换为JSON字符串
	jsonBytes, err := json.Marshal(l.queue)
	if err != nil {
		fmt.Println("JSON编码出错:", err)
		return
	}
	l.write(string(jsonBytes))

	// 重置缓冲区，准备接收新的日志消息
	l.count = 0
	l.queue = l.queue[:0]
	defer l.wg.Done()
}

func SetMode(mode string) {
	switch mode {
	case "text":
		OutFormat = 0
	case "json":
		OutFormat = 1
	default:
		log.Println("未知Mode")
	}
}

func file() *os.File {
	checkOutput() // 检查并创建日志输出目录

	// 创建日志文件，文件名由当前时间的日期格式和".log"后缀组成
	output, err := os.Create("./log/" + OutputPrefix + Mode[OutFormat])
	if err != nil { // 如果创建文件失败，记录错误信息
		log.Println(err)
	}
	return output
}
