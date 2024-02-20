package tinylog

import (
	"bufio"
	"log"
	"os"
)

// Write 方法用于将Logger缓冲区中的日志消息写入到文件，并刷新缓冲区。
func (l *Logger) Write() {
	// 创建一个带有缓冲的写入器，提高写入效率
	writer := bufio.NewWriter(l.Output)

	// 尝试将缓冲区的内容写入文件
	_, err := writer.WriteString(l.Buf.String())
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
		err = os.MkdirAll(logDir, 0755) // 创建目录并设置适当的权限
		if err != nil {
			log.Println("创建日志目录失败:", err)
		}
	}
}
