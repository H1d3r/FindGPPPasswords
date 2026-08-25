// Package logger gates the Manticore logger behind the --quiet flag.
//
// Manticore's logger has no quiet mode of its own, so the suppression lives
// here while the formatting and printing are delegated to the library.
package logger

import (
	manticore "github.com/TheManticoreProject/Manticore/logger"
)

var quiet bool

// SetQuiet enables or disables the suppression of all log output.
func SetQuiet(value bool) {
	quiet = value
}

func Info(message string) {
	if quiet {
		return
	}
	manticore.Info(message)
}

func Warn(message string) {
	if quiet {
		return
	}
	manticore.Warn(message)
}

func Debug(message string) {
	if quiet {
		return
	}
	manticore.Debug(message)
}
