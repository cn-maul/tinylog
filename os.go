package tinylog

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type jsonData struct {
	Data []LogMsg
}

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

func write(output io.Writer, input string) {
	buf := bufio.NewWriter(output)
	_, err := buf.WriteString(input)
	if err != nil {
		log.Println(err)
		return
	}

	err = buf.Flush()
	if err != nil {
		log.Println(err)
		return
	}
}

func open() *os.File {
	logDir := "./log"
	// 检查日志目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 如果目录不存在，创建它
		log.Println("日志文件夹不存在，正在创建log目录")
		// 创建目录并设置适当的权限
		if err = os.MkdirAll(logDir, 0o755); err != nil {
			log.Println("创建log目录失败:", err)
			return nil
		}
	}

	filePath := "./log/" + Filename + Mode[OutFormat]
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil { // 如果创建文件失败，记录错误信息
		log.Println(err)
	}
	return file
}

func (l *Logger) writeLog() {
	var buf strings.Builder
	for _, v := range l.queue {
		content := l.Prefix + Level[v.level] + v.msg + "\n"
		buf.Grow(len(content))
		buf.WriteString(content)
	}
	write(l.Output, buf.String())
	// 重置缓冲区，准备接收新的日志消息
	l.count = 0
	l.queue = l.queue[:0]
}

func (l *Logger) writeJSON() {
	data := jsonData{
		Data: l.queue,
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	jsonBytes, _ := json.Marshal(&data)
	write(l.Output, string(jsonBytes))
	// 重置缓冲区，准备接收新的日志消息
	l.count = 0
	l.queue = l.queue[:0]
}

// Flush 方法用于将Logger缓冲区中的所有日志消息写入到日志文件中。
func (l *Logger) Flush() {
	switch OutFormat {
	case 0:
		l.writeLog()
	case 1:
		l.writeJSON()
	default:
		log.Println("Flush时mode出错")
	}
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
