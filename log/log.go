package log

import (
	"log"
	"os"
)

var (
	lInfo    = log.New(os.Stdout, "[Info] ", log.Ldate|log.Ltime|log.Lshortfile)
	lWarning = log.New(os.Stdout, "[Warning] ", log.Ldate|log.Ltime|log.Lshortfile)
	lError   = log.New(os.Stderr, "[Error] ", log.Ldate|log.Ltime|log.Lshortfile)
	lFatal   = log.New(os.Stderr, "[Fatal] ", log.Ldate|log.Ltime|log.Lshortfile)
)

func Info(format string, v ...interface{}) {
	lInfo.Printf(format, v...)
}

func Warning(format string, v ...interface{}) {
	lWarning.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	lError.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	lFatal.Fatalf(format, v...)
}
