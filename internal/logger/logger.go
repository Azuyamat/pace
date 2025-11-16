package logger

import (
	"fmt"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorWhite  = "\033[97m"
)

type Logger struct {
	enabled bool
}

func New() *Logger {
	return &Logger{enabled: true}
}

func (l *Logger) SetEnabled(enabled bool) {
	l.enabled = enabled
}

func (l *Logger) timestamp() string {
	return colorGray + time.Now().Format("15:04:05") + colorReset
}

func (l *Logger) Info(format string, args ...interface{}) {
	if !l.enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s%s%s\n", l.timestamp(), colorCyan, msg, colorReset)
}

func (l *Logger) Success(format string, args ...interface{}) {
	if !l.enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s✓%s %s\n", l.timestamp(), colorGreen, colorReset, msg)
}

func (l *Logger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s✗%s %s\n", l.timestamp(), colorRed, colorReset, msg)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	if !l.enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s⚠%s %s\n", l.timestamp(), colorYellow, colorReset, msg)
}

func (l *Logger) Task(format string, args ...interface{}) {
	if !l.enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s▶%s %s\n", l.timestamp(), colorBlue, colorReset, msg)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s[DEBUG]%s %s\n", l.timestamp(), colorGray, colorReset, msg)
}

func (l *Logger) Print(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(msg)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *Logger) Println(args ...interface{}) {
	fmt.Println(args...)
}

var Default = New()

func Info(format string, args ...interface{}) {
	Default.Info(format, args...)
}

func Success(format string, args ...interface{}) {
	Default.Success(format, args...)
}

func Error(format string, args ...interface{}) {
	Default.Error(format, args...)
}

func Fatal(err error) error {
	Default.Error("%v", err)
	return err
}

func Warning(format string, args ...interface{}) {
	Default.Warning(format, args...)
}

func Task(format string, args ...interface{}) {
	Default.Task(format, args...)
}

func Debug(format string, args ...interface{}) {
	Default.Debug(format, args...)
}

func Print(format string, args ...interface{}) {
	Default.Print(format, args...)
}

func Printf(format string, args ...interface{}) {
	Default.Printf(format, args...)
}

func Println(args ...interface{}) {
	Default.Println(args...)
}
