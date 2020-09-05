package log

import (
	"bytes"
	"context"
	"fmt"
	"gitee.com/kelvins-io/common/env"
	"github.com/rs/zerolog"
	"os"
	"runtime"
	"time"
)

const (
	LEVEL_DEBUG = uint8(iota)
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR

	ZEROLOG_CALLER_SKIP = 3
	LOG_CALLER_SKIP     = 4
)

func init() {
	zerolog.TimestampFunc = func() time.Time {
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			return time.Now()
		}

		return time.Now().In(location)
	}
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.MessageFieldName = "msg"
}

type Fields = map[string]interface{}

func newLoggerContext() *LoggerContext {
	logger := &LoggerContext{}
	logger.loggerConf = LoggerConf{}
	return logger
}

type LoggerConf struct {
	LogFilePath string // Log文件全路径
	AppName     string
	Level       uint8
	Caller      bool // 是否需要显示Caller，默认false时不显示（注：获取caller会影响性能）
	HideTime    bool // 是否隐藏时间，默认显示
	ContextCall func(l *LoggerContext, ctx context.Context) LoggerContextIface
}

type LoggerContextIface interface {
	WithLogData(logData *LogData) LoggerContextIface
	WithCommonFields(commonFields Fields) LoggerContextIface
	WithFields(fields Fields) LoggerContextIface
	WithContext(ctx context.Context) LoggerContextIface
	Tag(tag string) LoggerContextIface
	RequestId(requestId string) LoggerContextIface
	WithCaller(skip int) LoggerContextIface

	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})

	GetDebugLog() *zerolog.Event
	GetInfoLog() *zerolog.Event
	GetWarnLog() *zerolog.Event
	GetErrorLog() *zerolog.Event

	GetLoggerConf() LoggerConf
}

type LoggerContext struct {
	logger *zerolog.Logger

	tag          string
	requestId    string
	commonFields []Fields
	fields       []Fields
	callerSkip   int

	loggerConf LoggerConf

	ctx context.Context
}

// 用于全局设置logger对象的公共输出字段
// 例如：logger = logger.WithCommonFields(logger.Fields{"app-name":"project" })
// {"level":"info",app-name":"project", "attach":{"process-id":25921,"speed":"33.771µs"},"msg":"for test"}
func (l *LoggerContext) WithCommonFields(commonFields Fields) LoggerContextIface {
	logger := l.clone()
	logger.commonFields = append(logger.commonFields, commonFields)
	return logger
}

// 用于输出自定义字段，同时又不影响整体的日志格式（将fields中指定的字段及值附加到attach中）
// {"level":"info","attach":{"process-id":25921,"speed":"33.771µs"},"msg":"for test"}
func (l *LoggerContext) WithFields(fields Fields) LoggerContextIface {
	logger := l.clone()
	logger.fields = append(logger.fields, fields)
	return logger
}

func (l *LoggerContext) WithContext(ctx context.Context) LoggerContextIface {
	logger := l.clone()
	logger.ctx = ctx
	if ctx != nil && logger.loggerConf.ContextCall != nil {
		logger = logger.loggerConf.ContextCall(logger, ctx).(*LoggerContext)
	}
	return logger
}

func (l *LoggerContext) Tag(tag string) LoggerContextIface {
	logger := l.clone()
	logger.tag = tag
	return logger
}

func (l *LoggerContext) RequestId(requestId string) LoggerContextIface {
	logger := l.clone()
	logger.requestId = requestId
	return logger
}

func (l *LoggerContext) WithCaller(skip int) LoggerContextIface {
	logger := l.clone()
	logger.callerSkip = skip
	return logger
}

func (l *LoggerContext) Debug(ctx context.Context, args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, args...)
	l.WithContext(ctx).(*LoggerContext).log(LEVEL_DEBUG, b.String())
}

func (l *LoggerContext) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.WithContext(ctx).(*LoggerContext).log(LEVEL_DEBUG, fmt.Sprintf(format, args...))
}

func (l *LoggerContext) Info(ctx context.Context, args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, args...)
	l.WithContext(ctx).(*LoggerContext).log(LEVEL_INFO, b.String())
}

func (l *LoggerContext) Infof(ctx context.Context, format string, args ...interface{}) {
	l.WithContext(ctx).(*LoggerContext).log(LEVEL_INFO, fmt.Sprintf(format, args...))
}

func (l *LoggerContext) Warn(ctx context.Context, args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, args...)
	l.WithContext(ctx).(*LoggerContext).log(LEVEL_WARN, b.String())
}

func (l *LoggerContext) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.WithContext(ctx).(*LoggerContext).log(LEVEL_WARN, fmt.Sprintf(format, args...))
}

func (l *LoggerContext) Error(ctx context.Context, args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(&b, args...)
	l.WithContext(ctx).(*LoggerContext).log(LEVEL_ERROR, b.String())
}

func (l *LoggerContext) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.WithContext(ctx).(*LoggerContext).log(LEVEL_ERROR, fmt.Sprintf(format, args...))
}

func (l *LoggerContext) GetDebugLog() *zerolog.Event {
	loggerEvent := l.logger.Debug()
	l.withCommon(loggerEvent, ZEROLOG_CALLER_SKIP)
	return loggerEvent
}

func (l *LoggerContext) GetInfoLog() *zerolog.Event {
	loggerEvent := l.logger.Info()
	l.withCommon(loggerEvent, ZEROLOG_CALLER_SKIP)
	return loggerEvent
}

func (l *LoggerContext) GetWarnLog() *zerolog.Event {
	loggerEvent := l.logger.Warn()
	l.withCommon(loggerEvent, ZEROLOG_CALLER_SKIP)
	return loggerEvent
}

func (l *LoggerContext) GetErrorLog() *zerolog.Event {
	loggerEvent := l.logger.Error()
	l.withCommon(loggerEvent, ZEROLOG_CALLER_SKIP)
	return loggerEvent
}

func (l *LoggerContext) GetLoggerConf() LoggerConf {
	return l.loggerConf
}

func (l *LoggerContext) getCallerSkip(defaultCallerSkip int) int {
	if l.callerSkip > 0 {
		return l.callerSkip
	}
	return defaultCallerSkip
}

func (l *LoggerContext) log(level uint8, msg string) {
	var zeroEvent *zerolog.Event
	switch level {
	case LEVEL_DEBUG:
		zeroEvent = l.logger.Debug()
	case LEVEL_INFO:
		zeroEvent = l.logger.Info()
	case LEVEL_WARN:
		zeroEvent = l.logger.Warn()
	case LEVEL_ERROR:
		if env.IsDevMode() {
			l := zerolog.New(os.Stderr)
			l.Error().Caller().Timestamp().Msg(msg)
		}
		zeroEvent = l.logger.Error()
	}
	l.withCommon(zeroEvent, LOG_CALLER_SKIP)
	zeroEvent.Msg(msg)
}

func (l *LoggerContext) generateCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d %s", file, line, f.Name())
}

func (l *LoggerContext) clone() *LoggerContext {
	newLog := *l
	return &newLog
}

func (l *LoggerContext) withCommon(loggerEvent *zerolog.Event, skip int) {
	for _, field := range l.commonFields {
		loggerEvent.Fields(field)
	}

	if !l.loggerConf.HideTime {
		loggerEvent.Timestamp()
	}
	if l.tag != "" {
		loggerEvent.Str("tag", l.tag)
	}
	if l.requestId != "" {
		loggerEvent.Str("request-id", l.requestId)
	}
	if l.loggerConf.Caller {
		loggerEvent.Str("caller", l.generateCallerInfo(l.getCallerSkip(skip)))
	}
	dict := zerolog.Dict()
	for _, field := range l.fields {
		dict.Fields(field)
	}
	if len(l.fields) > 0 {
		loggerEvent.Dict("attach", dict)
	}
}
