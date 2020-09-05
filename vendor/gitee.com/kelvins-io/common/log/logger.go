package log

import (
	"errors"

	"github.com/rs/zerolog"
)

// 创建日志对象
func CreateLogger(logConf LoggerConf) (LoggerContextIface, error) {
	if logConf.LogFilePath == "" {
		return nil, errors.New("Log Path is empty")
	}
	trackLogger := newLoggerContext()
	trackLogger.loggerConf = logConf
	fileWriter, err := newFileWriter(logConf.LogFilePath)
	if err != nil {
		return nil, err
	}
	ioLogger := zerolog.New(fileWriter)
	ioLogger.Level(zerolog.Level(trackLogger.loggerConf.Level))

	trackLogger.logger = &ioLogger
	return trackLogger, nil
}
