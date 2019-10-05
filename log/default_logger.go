package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type defaultLogger struct {
	level int
	*log.Logger
}

func (l *defaultLogger) SetLevel(level int) {
	l.level = level
}

func (l *defaultLogger) SprintForOutput(arr ...interface{}) string {
	var ret strings.Builder
	for k, _ := range arr {
		ret.WriteString(fmt.Sprint(arr[k]))
		ret.WriteString(" ")
	}
	return ret.String()
}

func (l *defaultLogger) Debug(v ...interface{}) {
	if l.level > LDebug {
		return
	}
	l.Output(calldepth, header("DEBUG", l.SprintForOutput(v...)))
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	if l.level > LDebug {
		return
	}
	l.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Info(v ...interface{}) {
	if l.level > LInfo {
		return
	}
	// l.Output(calldepth, header(color.GreenString("INFO "), l.SprintForOutput(v...)))
	l.Output(calldepth, header("INFO ", l.SprintForOutput(v...)))
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	if l.level > LInfo {
		return
	}
	// l.Output(calldepth, header(color.GreenString("INFO "), fmt.Sprintf(format, v...)))
	l.Output(calldepth, header("INFO ", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Warn(v ...interface{}) {
	if l.level > LWarn {
		return
	}
	// l.Output(calldepth, header(color.YellowString("WARN "), l.SprintForOutput(v...)))
	l.Output(calldepth, header("WARN ", l.SprintForOutput(v...)))
}

func (l *defaultLogger) Warnf(format string, v ...interface{}) {
	if l.level > LWarn {
		return
	}
	// l.Output(calldepth, header(color.YellowString("WARN "), fmt.Sprintf(format, v...)))
	l.Output(calldepth, header("WARN ", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Error(v ...interface{}) {
	if l.level > LError {
		return
	}
	// l.Output(calldepth, header(color.RedString("ERROR"), l.SprintForOutput(v...)))
	l.Output(calldepth, header("ERROR", l.SprintForOutput(v...)))
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	if l.level > LError {
		return
	}
	// l.Output(calldepth, header(color.RedString("ERROR"), fmt.Sprintf(format, v...)))
	l.Output(calldepth, header("ERROR", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Fatal(v ...interface{}) {
	if l.level > LFatal {
		return
	}
	// l.Output(calldepth, header(color.MagentaString("FATAL"), l.SprintForOutput(v...)))
	l.Output(calldepth, header("FATAL", l.SprintForOutput(v...)))
	os.Exit(1)
}

func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	if l.level > LFatal {
		return
	}
	// l.Output(calldepth, header(color.MagentaString("FATAL"), fmt.Sprintf(format, v...)))
	l.Output(calldepth, header("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *defaultLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v)
}

func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}

func header(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}
