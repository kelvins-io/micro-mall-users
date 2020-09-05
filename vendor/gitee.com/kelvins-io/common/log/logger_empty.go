package log

import (
	"context"
	"github.com/rs/zerolog"
)

type EmptyLoggerContext struct {
}

func (*EmptyLoggerContext) WithLogData(logData *LogData) LoggerContextIface {
	return &EmptyLoggerContext{}
}

func (*EmptyLoggerContext) WithCommonFields(fields Fields) LoggerContextIface {
	return &EmptyLoggerContext{}
}
func (*EmptyLoggerContext) WithFields(fields Fields) LoggerContextIface {
	return &EmptyLoggerContext{}
}
func (*EmptyLoggerContext) WithContext(ctx context.Context) LoggerContextIface {
	return &EmptyLoggerContext{}
}

func (*EmptyLoggerContext) Tag(tag string) LoggerContextIface {
	return &EmptyLoggerContext{}
}

func (*EmptyLoggerContext) RequestId(requestId string) LoggerContextIface {
	return &EmptyLoggerContext{}
}

func (*EmptyLoggerContext) WithCaller(skip int) LoggerContextIface {
	return &EmptyLoggerContext{}
}

func (e *EmptyLoggerContext) Debug(ctx context.Context, msg ...interface{})                  {}
func (e *EmptyLoggerContext) Debugf(ctx context.Context, format string, args ...interface{}) {}
func (e *EmptyLoggerContext) Info(ctx context.Context, msg ...interface{})                   {}
func (e *EmptyLoggerContext) Infof(ctx context.Context, format string, args ...interface{})  {}
func (e *EmptyLoggerContext) Warn(ctx context.Context, msg ...interface{})                   {}
func (e *EmptyLoggerContext) Warnf(ctx context.Context, format string, args ...interface{})  {}
func (e *EmptyLoggerContext) Error(ctx context.Context, msg ...interface{})                  {}
func (e *EmptyLoggerContext) Errorf(ctx context.Context, format string, args ...interface{}) {}

func (e *EmptyLoggerContext) GetDebugLog() *zerolog.Event {
	event := &zerolog.Event{}
	event.Discard()
	return event
}

func (e *EmptyLoggerContext) GetInfoLog() *zerolog.Event {
	event := &zerolog.Event{}
	event.Discard()
	return event
}

func (e *EmptyLoggerContext) GetWarnLog() *zerolog.Event {
	event := &zerolog.Event{}
	event.Discard()
	return event
}

func (e *EmptyLoggerContext) GetErrorLog() *zerolog.Event {
	event := &zerolog.Event{}
	event.Discard()
	return event
}

func (e *EmptyLoggerContext) GetLoggerConf() LoggerConf {
	return LoggerConf{}
}
