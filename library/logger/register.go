package logger

type LoggerOption func(*LoggerOptions)

var (
	DefaultOptions = LoggerOptions{
		LogDirName:               "./logs",
		AllowLogLevel:            LevelW,
		AllowAccessLog:           true,
		AllowInterfaceAvgTimeLog: true,
		FlushInterval:            3,
	}
)
