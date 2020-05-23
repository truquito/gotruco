package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

// LogFile log file
type LogFile struct {
	path string
}

// Log str to logFile
func (lf LogFile) Write(str string) {
	f, err := os.OpenFile(lf.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(str)
}

func newLogFile(logPath string) LogFile {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	logFile := logPath + timestamp + ".log"

	return LogFile{path: logFile}
}
