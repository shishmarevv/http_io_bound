package errlog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Report struct {
	Info  *log.Logger
	Error *log.Logger
	File  *os.File
}

func ReportFile() (*os.File, error) {
	logsDir := os.Getenv("LOG_DIR")
	if logsDir == "" {
		logsDir = "logs"
	}

	var path string
	path = filepath.Join(logsDir, fmt.Sprintf("%s.log", time.Now().Format("02-01-2006")))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func NewReport() (*Report, error) {
	file, err := ReportFile()
	if err != nil {
		return nil, err
	}
	infolog := log.New(file, "[INFO] ", log.Ltime)
	errorlog := log.New(file, "[ERROR] ", log.Ltime|log.Lshortfile)
	return &Report{
		Info:  infolog,
		Error: errorlog,
		File:  file,
	}, nil
}
