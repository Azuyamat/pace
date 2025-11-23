package logger

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	LevelInfo LogLevel = iota
	LevelWarning
	LevelError
	LevelDebug
)

type Logger struct {
	enabled bool
	level   LogLevel
}

func New() *Logger {
	if os.Getenv("PACE_DEBUG") == "true" {
		return &Logger{enabled: true, level: LevelDebug}
	}
	return &Logger{enabled: true, level: LevelInfo}
}

func (l *Logger) SetEnabled(enabled bool) {
	l.enabled = enabled
}

func (l *Logger) timestamp() string {
	background := ColorBlue.Dark().Background()
	return background.Wrap(" " + time.Now().Format("15:04:05") + " ")
}

func (l *Logger) Info(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	foregroundColor := ColorWhite
	icon := ColorBlue.Dim().Wrap("[INFO]")
	coloredMessage := foregroundColor.Wrap(msg)
	fmt.Printf("%s %s %s\n", l.timestamp(), icon, coloredMessage)
}

func (l *Logger) Success(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s✓%s %s\n", l.timestamp(), ColorGreen, ColorReset, msg)
}

func (l *Logger) Error(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelError {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s✗%s %s\n", l.timestamp(), ColorRed, ColorReset, msg)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelWarning {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s⚠%s %s\n", l.timestamp(), ColorYellow, ColorReset, msg)
}

func (l *Logger) Task(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s▶%s %s\n", l.timestamp(), ColorBlue, ColorReset, msg)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.enabled || l.level != LevelDebug {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s[DEBUG]%s %s\n", l.timestamp(), ColorGray, ColorReset, msg)
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
