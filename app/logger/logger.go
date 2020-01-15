package logger

import (
	"fmt"
	"log"
	"os"
)

type LEVEL int

const (
	debug LEVEL = iota
	info
	warning
	error
	fatal
)

var (
	logPrefix = ""

	levelFlags = []string{"DEBG", "INFO", "WARN", "ERRO", "FATL"}

	logInstance *log.Logger

	logfile *os.File
)

// Debug ...
func Debug(v ...interface{}) {
	setPrefix(debug)
	println(logInstance, v)

}

// Info ...
func Info(v ...interface{}) {
	setPrefix(info)
	println(logInstance, v)
}

// Warn ...
func Warn(v ...interface{}) {
	setPrefix(warning)
	println(logInstance, v)
}

// Error Warn
func Error(v ...interface{}) {
	setPrefix(error)
	println(logInstance, v)
}

// Fatal ...
func Fatal(v ...interface{}) {
	setPrefix(fatal)
	fatalln(logInstance, v)
}

func setPrefix(level LEVEL) {
	logPrefix = fmt.Sprintf("[%s] ", levelFlags[level])
	logInstance.SetPrefix(logPrefix)
}

// Println ..
func println(l *log.Logger, v ...interface{}) {
	if l != nil {
		l.Output(3, fmt.Sprintln(v...))
	}

}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func fatalln(l *log.Logger, v ...interface{}) {
	if l != nil {
		l.Output(3, fmt.Sprintln(v...))
		os.Exit(1)
	}
}
