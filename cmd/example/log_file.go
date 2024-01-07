package main

import (
	"fmt"
	"log"
	"os"
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
	// Get current date and time in YYYY-MM-DD HH:MM:SS format
	now := time.Now().UTC()
	formattedTime := fmt.Sprintf("%d-%02d-%02d_%02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())

	// Create log file name with formatted date and time
	logFile := logPath + formattedTime + ".log"

	return LogFile{path: logFile}
}
