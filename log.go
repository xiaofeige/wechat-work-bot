package wechat_work_bot

import "fmt"

// Logger interface API for log.Logger
var (
	DefaultLogger = &defaultLogger{}
)

type defaultLogger struct {
}

func (l *defaultLogger)Printf(format string, a ...interface{})  {
	fmt.Println(fmt.Sprintf(format, a))
}

type Logger interface {
	Printf(string, ...interface{})
}

