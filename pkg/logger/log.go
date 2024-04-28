package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) LogInfo(msg string, fileName *string) {
	l.log(INFO, msg, fileName)
}

func (l *Logger) LogWarning(msg string, fileName *string) {
	l.log(WARN, msg, fileName)
}

func (l *Logger) LogError(msg string, fileName *string) {
	l.log(ERROR, msg, fileName)
}

func (l *Logger) LogDebug(msg string, fileName *string) {
	l.log(DEBUG, msg, fileName)
}

func (l *Logger) log(level string, msg string, fileName *string) {
	_, file, line, _ := runtime.Caller(2)
	logMsg := fmt.Sprintf("%s | %v | %v | %v:%v | %v\n", time.Now().Format(time.RFC3339), level, "", file, line, msg)

	if fileName != nil {

		f, err := os.OpenFile(getPath(*fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
			return
		}
		defer f.Close()

		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime)) // remove prefix time

		log.SetOutput(f)
		log.Println(logMsg)
	}

	fmt.Println(logMsg)
}

func getPath(fileName string) string {
	path := ""
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path = dir + "\\logs\\"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	return path + "\\" + fileName + ".log"
}
