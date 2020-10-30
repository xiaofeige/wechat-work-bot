package wechat_work_bot

import "fmt"

// Logger interface API for log.Logger
var (
	DefaultLogger Logger
)

func init() {
	DefaultLogger = LoggerFunc(func(msg string, args ...interface{}) {
		fmt.Printf(msg, args...)
	})
}

type Logger interface {
	Printf(string, ...interface{})
}

// LoggerFunc is a bridge between Logger and any third party logger
// Usage:
//   l := NewLogger() // some logger
//   r := wechat_work_bot.NewRobot(wechat_work_bot.RobotConfig{
//     Debugger:	wechat_work_bot.LoggerFunc(l.Infof),
//     ErrLogger: 	wechat_work_bot.LoggerFunc(l.Errorf),
//   })
type LoggerFunc func(string, ...interface{})

func (f LoggerFunc) Printf(msg string, args ...interface{}) { f(msg, args...) }

// replace global default logger with user provide logger
func SetLogger(logger Logger) {
	DefaultLogger = logger
}
