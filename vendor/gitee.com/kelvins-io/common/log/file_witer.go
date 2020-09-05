package log

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	getLogFilePath = func(logPath string) string {
		logFileFullPath := logPath + time.Now().Format(".2006-01-02")
		return logFileFullPath
	}
)

type fileWriter struct {
	file            *os.File
	logFileFullPath string
	logFilePath     string
	mu              sync.Mutex
}

func (f *fileWriter) Write(b []byte) (int, error) {
	f.mu.Lock()
	logFileFullPath := getLogFilePath(f.logFilePath)
	if logFileFullPath != f.logFileFullPath {
		file, err := os.OpenFile(logFileFullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			f.mu.Unlock()
			return 0, err
		}
		_ = f.file.Close()
		f.file = file
		f.logFileFullPath = logFileFullPath
	}
	f.mu.Unlock()

	// os.File Write operation is using write() syscall which is supposed to be thread-safe on POSIX systems.
	n, err := f.file.Write(b)
	return n, err
}

func newFileWriter(filePath string) (*fileWriter, error) {
	logFileFullPath := getLogFilePath(filePath)

	dirName := filepath.Dir(logFileFullPath)
	if err := os.MkdirAll(dirName, 0775); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logFileFullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &fileWriter{
		file:            file,
		logFileFullPath: logFileFullPath,
		logFilePath:     filePath,
	}, nil
}
