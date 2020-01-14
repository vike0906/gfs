package logger

import (
	"bufio"
	"gfs/app/util"
	"log"
	"os"
	"time"
)

const bufSize = 2 * 1024

var (
	logWriter *bufio.Writer
	logPath   string
)

func GetLogWriter() *bufio.Writer {
	return logWriter
}

//begin when app start
//and then in the 00:00:00 every day
func createLogFileTask() {
	go func() {
		for {
			logFileName := util.NewDateHelper().FormatDate(time.Now()) + ".log"
			file, err := os.OpenFile(logPath+logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Println(err.Error())
				logWriter = nil
			} else {
				logWriter = bufio.NewWriterSize(file, bufSize)
				//log.Println(LogWriter)
			}
			//lock
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			_ = logWriter.Flush()
		}
	}()
}

//check is logs folder exist, mkdir if not
func init() {
	logFolder := "/logs/"
	if logFolderPath, err := util.PathAdaptive(logFolder); err != nil {
		log.Fatal("log model init failed")
	} else {
		logPath = logFolderPath
		if _, err := os.Stat(logFolderPath); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(logFolderPath, os.ModePerm); err != nil {
					log.Fatal(err.Error())
				}
			}
		}
	}
	createLogFileTask()
}
