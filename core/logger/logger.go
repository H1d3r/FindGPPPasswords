package logger

import (
	"fmt"
	"time"
)

var quiet bool

func SetQuiet(value bool) {
	quiet = value
}

func Info(message string) {
	Dateprintf("INFO: %s\n", message)
}

func Warn(message string) {
	Dateprintf("WARN: %s\n", message)
}

func Debug(message string) {
	Dateprintf("DEBUG: %s\n", message)
}

func Dateprintf(format string, message ...any) {
	if quiet {
		return
	}
	currentTime := time.Now().Format("2006-01-02 15h04m05s")
	format = fmt.Sprintf("[%s] %s", currentTime, format)
	fmt.Printf(format, message...)
}
