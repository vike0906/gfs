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
	Println(logInstance, v)

}

// Info ...
func Info(v ...interface{}) {
	setPrefix(info)
	Println(logInstance, v)
}

// Warn ...
func Warn(v ...interface{}) {
	setPrefix(warning)
	Println(logInstance, v)
}

// Error Warn
func Error(v ...interface{}) {
	setPrefix(error)
	Println(logInstance, v)
}

// Fatal ...
func Fatal(v ...interface{}) {
	setPrefix(fatal)
	Fatalln(logInstance, v)
}

func setPrefix(level LEVEL) {
	logPrefix = fmt.Sprintf("[%s] ", levelFlags[level])
	logInstance.SetPrefix(logPrefix)
}

// Println ..
func Println(l *log.Logger, v ...interface{}) {
	if l != nil {
		l.Output(3, fmt.Sprintln(v...))
	}

}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func Fatalln(l *log.Logger, v ...interface{}) {
	if l != nil {
		l.Output(3, fmt.Sprintln(v...))
		os.Exit(1)
	}
}
