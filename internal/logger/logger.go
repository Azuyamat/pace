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
	timeStr := time.Now().Format("15:04:05")
	return ColorGray.Wrap("[" + timeStr + "]")
}

func (l *Logger) badge(text string, color Color) string {
	bg := color.Dark().Background()
	return bg.Wrap(" " + text + " ")
}

func (l *Logger) Info(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	badge := l.badge("INFO", ColorBlue)
	icon := ColorCyan.Wrap("◆")
	coloredMsg := ColorWhite.Wrap(msg)
	fmt.Printf("%s %s %s %s\n", l.timestamp(), badge, icon, coloredMsg)
}

func (l *Logger) Success(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	badge := l.badge("DONE", ColorGreen)
	icon := ColorGreen.Bright().Wrap("✓")
	coloredMsg := ColorWhite.Wrap(msg)
	fmt.Printf("%s %s %s %s\n", l.timestamp(), badge, icon, coloredMsg)
}

func (l *Logger) Error(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelError {
		return
	}
	msg := fmt.Sprintf(format, args...)
	badge := l.badge("ERROR", ColorRed)
	icon := ColorRed.Bright().Wrap("✗")
	coloredMsg := ColorWhite.Wrap(msg)
	fmt.Printf("%s %s %s %s\n", l.timestamp(), badge, icon, coloredMsg)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelWarning {
		return
	}
	msg := fmt.Sprintf(format, args...)
	badge := l.badge("WARN", ColorYellow)
	icon := ColorYellow.Bright().Wrap("⚠")
	coloredMsg := ColorWhite.Wrap(msg)
	fmt.Printf("%s %s %s %s\n", l.timestamp(), badge, icon, coloredMsg)
}

func (l *Logger) Task(format string, args ...interface{}) {
	if !l.enabled || l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	badge := l.badge("TASK", ColorPurple)
	icon := ColorPurple.Bright().Wrap("▶")
	coloredMsg := ColorWhite.Wrap(msg)
	fmt.Printf("%s %s %s %s\n", l.timestamp(), badge, icon, coloredMsg)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.enabled || l.level != LevelDebug {
		return
	}
	msg := fmt.Sprintf(format, args...)
	badge := l.badge("DEBUG", ColorGray)
	icon := ColorGray.Wrap("●")
	coloredMsg := ColorGray.Wrap(msg)
	fmt.Printf("%s %s %s %s\n", l.timestamp(), badge, icon, coloredMsg)
}

func (l *Logger) TaskOutput(taskName string, format string, args ...interface{}) {
	if !l.enabled || l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	taskBadge := ColorCyan.Wrap("[" + taskName + "]")
	fmt.Printf("%s %s %s\n", l.timestamp(), taskBadge, ColorWhite.Wrap(msg))
}

func (l *Logger) TaskError(taskName string, format string, args ...interface{}) {
	if !l.enabled || l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	taskBadge := ColorRed.Wrap("[" + taskName + "]")
	fmt.Printf("%s %s %s\n", l.timestamp(), taskBadge, ColorRed.Bright().Wrap(msg))
}

func (l *Logger) Prompt(format string, args ...interface{}) (string, error) {
	if !l.enabled || l.level > LevelInfo {
		return "", fmt.Errorf("logger is disabled")
	}
	msg := fmt.Sprintf(format, args...)
	promptBadge := l.badge("PROMPT", ColorCyan)
	fmt.Printf("%s %s %s ", l.timestamp(), promptBadge, ColorWhite.Wrap(msg))
	var input string
	fmt.Scanln(&input)
	return input, nil
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

func Prompt(format string, args ...interface{}) (string, error) {
	return Default.Prompt(format, args...)
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
