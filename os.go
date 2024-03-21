package tlog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Suffix 日志文件后缀名
var Suffix = ".log"

// OutputDir 日志文件输出目录
var OutputDir = "log/"

// 写入日志
func write() {
	switch Suffix {
	case ".json":
		writeJSON()
	case ".log":
		writePlain()
	default:
		log.Println("后缀不存在")
	}
}

// writeJSON 写入JSON格式日志
func writeJSON() {
	type JSON struct {
		Log []Log
	}
	bytes, err := json.Marshal(&JSON{Log: List})
	if err != nil {
		log.Println(err)
		return
	}
	fileName := time.Now().Format("2006-01-02 15-04-05")
	os.WriteFile(OutputDir+fileName+Suffix, bytes, 0777)
	List = List[:0]
}

// writePlain 写入纯文本格式日志
func writePlain() {
	var buf strings.Builder
	for _, v := range List {
		str := fmt.Sprintf("%s [%d] %s\n", v.Time, v.Level, v.Message)
		buf.Grow(len(str))
		buf.WriteString(str)
	}
	fileName := time.Now().Format("2006-01-02 15-04-05")
	os.WriteFile(OutputDir+fileName+Suffix, []byte(buf.String()), 0777)
	List = List[:0]
}
