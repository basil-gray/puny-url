package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	isSilent    bool
)

func Init(silent bool) {
	isSilent = silent
	if !isSilent {
		infoLogger = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime)
	}

	warnLogger = log.New(os.Stdout, "Warn: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(msg string) {
	if !isSilent {
		infoLogger.Println(msg)
	}
}

func Warn(msg string) {
	warnLogger.Println(msg)
}

func Error(msg string) {
	errorLogger.Println(msg)
}
