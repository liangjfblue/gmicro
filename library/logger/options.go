package logger

import "time"

type LoggerOptions struct {
	LogDirName               string
	AllowLogLevel            int32
	AllowAccessLog           bool
	AllowInterfaceAvgTimeLog bool
	FlushInterval            time.Duration
}

func newOptions(opts ...LoggerOption) LoggerOptions {
	opt := DefaultOptions

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func LogDirName(logDirName string) LoggerOption {
	return func(o *LoggerOptions) {
		o.LogDirName = logDirName
	}
}

func AllowLogLevel(allowLogLevel int32) LoggerOption {
	return func(o *LoggerOptions) {
		o.AllowLogLevel = allowLogLevel
	}
}

func AllowAccessLog(allowAccessLog bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.AllowAccessLog = allowAccessLog
	}
}

func AllowInterfaceAvgTimeLog(allowInterfaceAvgTimeLog bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.AllowInterfaceAvgTimeLog = allowInterfaceAvgTimeLog
	}
}

func FlushInterval(flushInterval time.Duration) LoggerOption {
	return func(o *LoggerOptions) {
		o.FlushInterval = flushInterval
	}
}
