package tinylog

import (
	"bufio"
	"log"
	"os"
)

func (l *Logger) Write() {
	//检查日志输出文件夹是否存在
	checkOutput()
	// 创建一个带有缓冲区的写入器
	writer := bufio.NewWriter(l.Output)

	// 写入字符串
	_, err := writer.WriteString(l.Buf.String())
	if err != nil {
		log.Println("Bufio写时出错:", err)
		return
	}

	// 刷新缓冲区，确保所有数据写入文件
	err = writer.Flush()
	if err != nil {
		log.Println("Flush时出错", err)
		return
	}
}

func checkOutput() {
	logDir := "./log"
	// 检查log文件夹是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 如果不存在，创建log文件夹
		log.Println("日志文件夹不存在，创建log目录")
		err = os.MkdirAll(logDir, 0755) // 创建目录并设置权限
	}
}
