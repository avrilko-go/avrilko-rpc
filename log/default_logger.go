package log

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

type defaultLogger struct {
	*log.Logger
}

func (d *defaultLogger) Debug(v ...interface{}) {
	d.Output(CallDepth, wrapperMsg("[DEBUG]", fmt.Sprint(v...)))
}

func (d *defaultLogger) DebugF(format string, v ...interface{}) {
	d.Output(CallDepth, wrapperMsg("[DEBUG]", fmt.Sprintf(format, v)))
}

func (d *defaultLogger) Info(v ...interface{}) {
	d.Output(CallDepth, wrapperMsg(color.GreenString("[INFO]"), fmt.Sprint(v...)))
}

func (d *defaultLogger) InfoF(format string, v ...interface{}) {
	d.Output(CallDepth, wrapperMsg(color.GreenString("[INFO]"), fmt.Sprintf(format, v)))
}

func (d *defaultLogger) Warn(v ...interface{}) {
	d.Output(CallDepth, wrapperMsg(color.YellowString("[WARN]"), fmt.Sprint(v...)))
}

func (d *defaultLogger) WarnF(format string, v ...interface{}) {
	d.Output(CallDepth, wrapperMsg(color.YellowString("[WARN]"), fmt.Sprintf(format, v)))
}

func (d *defaultLogger) Error(v ...interface{}) {
	d.Output(CallDepth, wrapperMsg(color.RedString("[ERROR]"), fmt.Sprint(v...)))
}

func (d *defaultLogger) ErrorF(format string, v ...interface{}) {
	d.Output(CallDepth, wrapperMsg(color.RedString("[ERROR]"), fmt.Sprintf(format, v)))
}

func (d *defaultLogger) Fatal(v ...interface{}) {
	d.Output(CallDepth, wrapperMsg(color.MagentaString("[FATAL]"), fmt.Sprint(v...)))
}
func (d *defaultLogger) FatalF(format string, v ...interface{}) {
	d.Output(CallDepth, wrapperMsg(color.MagentaString("[FATAL]"), fmt.Sprintf(format, v)))
}

func (d *defaultLogger) Panic(v ...interface{}) {
	d.Logger.Panic(v...)
}
func (d *defaultLogger) PanicF(format string, v ...interface{}) {
	d.Logger.Panicf(format, v...)
}

func (d *defaultLogger) Handle(v ...interface{}) {
	d.Error(v)
}

func wrapperMsg(level, msg string) string {
	return fmt.Sprintf("%s: %s", level, msg)
}
