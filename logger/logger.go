package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogSeverity int

const (
	DEBUG LogSeverity = iota
	ERROR LogSeverity = iota
	INFO  LogSeverity = iota
)

func GetSeverityString(logType LogSeverity) string {
	switch logType {
	case DEBUG:
		return "[DEBUG]"
	case ERROR:
		return "[ERROR]"
	case INFO:
		return "[INFO]"
	default:
		return "[INFO]"
	}
}

func write(logType LogSeverity, message string) {
	f, err := os.OpenFile("alfredbot.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.Printf("%s: %s\n", GetSeverityString(logType), message)
}

func WriteError(message string, err error) {
	write(ERROR, fmt.Sprintf("%s %s", message, err))
}

func WriteInfo(message string) {
	write(INFO, message)
}
