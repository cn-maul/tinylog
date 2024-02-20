# Tinylog

## 安装

```bash
go get github.com/cn-maul/tinylog
```

## 使用

```go
func (l *Logger) Println(logLevel, msg string) 
//logLevel为日志等级，分log和err
//msg为日志具体信息
```

例子

```go
var Log = tinylog.New()
defer Log.Output.Close()

msg := "msg message"
err := "err message"

Log.Println("log", msg)
Log.Println("err", err)

Log.Write()

//Output
2024/02/20 21:27:02 [Msg]msg message
2024/02/20 21:27:02 [Err]err message
	
```