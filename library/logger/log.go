package logger

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/liangjfblue/gmicro/library/logger/glog"
)

const (
	SevDebug                = "debug"
	SevInfo                 = "info"
	SevWarn                 = "warning"
	SevError                = "error"
	SevAccess               = "access"
	SevInterfaceAvgDuration = "iavgd"
	SevFatal                = "fatal"

	LevelD = 1 // debug
	LevelI = 2 // info
	LevelW = 3 // warning
	LevelE = 4 // error
	LevelA = 5 // access
)

type Logger struct {
	logger *glog.Logger

	LoggerOptions LoggerOptions
}

func NewLogger(opts ...LoggerOption) *Logger {
	return &Logger{
		LoggerOptions: newOptions(opts...),
	}
}

func (l *Logger) Init() {
	if _, err := os.Stat(l.LoggerOptions.LogDirName); os.IsNotExist(err) {
		if err := os.Mkdir(l.LoggerOptions.LogDirName, os.ModePerm); err != nil {
			panic(err)
		}
	}

	l.logger = glog.NewLogger().
		LogDir(l.LoggerOptions.LogDirName).
		EnableLogHeader(true).
		EnableLogLink(false).
		FlushInterval(l.LoggerOptions.FlushInterval).
		HeaderFormat(func(buf *bytes.Buffer, l glog.Severity, ts time.Time, pid int, file string, line int) {
			switch l {
			case glog.InfoLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][INFO]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case glog.DebugLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][DEBUG]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case glog.WarnLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][WARN]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case glog.ErrorLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][ERROR]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case glog.AccessLog:
				_, _ = fmt.Fprintf(buf, "[%s]\t", ts.Format("2006-01-02 15:04:05"))
			case glog.InterfaceAvgDurationLog:
				_, _ = fmt.Fprintf(buf, "[%s]\t", ts.Format("2006-01-02 15:04:05"))
			}
		}).
		FileNameFormat(fileNameFormatFunc).
		Init()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		_, _ = os.Stdout.WriteString("flushing log ...\n")
		l.logger.Flush()
		_, _ = os.Stdout.WriteString("flush log done\n")
		signal.Reset(os.Interrupt)

		proc, err := os.FindProcess(syscall.Getpid())
		if err != nil {
			panic(err)
		}

		err = proc.Signal(os.Interrupt)
		if err != nil {
			panic(err)
		}

	}()

}

func logTag(severityLevel string) string {
	tag := SevInfo
	switch severityLevel {
	case glog.SevDebug:
		tag = SevDebug
	case glog.SevInfo:
		tag = SevInfo
	case glog.SevWarn:
		tag = SevWarn
	case glog.SevError:
		tag = SevError
	case glog.SevAccess:
		tag = SevAccess
	case glog.SevFatal:
		tag = SevFatal
	case glog.SevInterfaceAvgDuration:
		tag = SevInterfaceAvgDuration
	}

	return tag
}

func fileNameFormatFunc(severityLevel string, ts time.Time) string {
	var filename string
	tag := logTag(severityLevel)
	filename = fmt.Sprintf("%s.log.%04d-%02d-%02d",
		tag,
		ts.Year(),
		ts.Month(),
		ts.Day())
	return filename
}

func (l *Logger) SetLogLevel(level int) {
	atomic.StoreInt32(&l.LoggerOptions.AllowLogLevel, int32(level))
}

func (l *Logger) GetLogLevel() int {
	return int(atomic.LoadInt32(&l.LoggerOptions.AllowLogLevel))
}

func (l *Logger) Info(format string, args ...interface{}) {
	if atomic.LoadInt32(&l.LoggerOptions.AllowLogLevel) <= LevelI {
		l.logger.InfoDepth(1, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if atomic.LoadInt32(&l.LoggerOptions.AllowLogLevel) <= LevelW {
		l.logger.WarningDepth(1, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if atomic.LoadInt32(&l.LoggerOptions.AllowLogLevel) <= LevelD {
		l.logger.DebugDepth(1, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if atomic.LoadInt32(&l.LoggerOptions.AllowLogLevel) <= LevelE {
		l.logger.ErrorDepth(1, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Access(format string, args ...interface{}) {
	if l.LoggerOptions.AllowAccessLog {
		l.logger.AccessDepth(1, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) InterfaceAvgDuration(format string, args ...interface{}) {
	if l.LoggerOptions.AllowInterfaceAvgTimeLog {
		l.logger.InterfaceAvgDurationDepth(1, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) FlushLog() {
	l.logger.Flush()
}
