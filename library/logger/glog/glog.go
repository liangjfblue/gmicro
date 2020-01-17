// Go support for leveled logs, analogous to https://code.google.com/p/google-glog/
//
// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package glog implements logging analogous to the Google-internal C++ INFO/ERROR/V setup.
// It provides functions Info, Warning, Error, Fatal, plus formatting variants such as
// Infof.
//
// Basic examples:
//
//	glog.Info("Prepare to repel boarders")
//
//	glog.Fatalf("Initialization failed: %s", err)
//
//
// Log output is buffered and written periodically using Flush. Programs
// should call Flush before exiting to guarantee all log output is written.
//
// By default, all log statements write to files in a temporary directory.
// This package provides several flags that modify this behavior.
// As a result, flag.Parse must be called before any logging is done.
//
//	-logtostderr=false
//		Logs are written to standard error instead of to files.
//	-alsologtostderr=false
//		Logs are written to standard error as well as to files.
//	-sevthreshold=ERROR
//		Log events at or above this severity are outputed.
//	-log_dir=""
//		Log files will be written to this directory instead of the
//		default temporary directory.
//
//	Other flags provide aids to debugging.
//
//	-log_backtrace_at=""
//		When set to a file and line number holding a logging statement,
//		such as
//			-log_backtrace_at=gopherflakes.go:234
//		a stack trace will be written to the Info log whenever execution
//		hits that statement. (Unlike with -vmodule, the ".go" must be
//		present.)
package glog

import (
	"bytes"
	"fmt"
	stdLog "log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func init() {

	// init std logger
	logging.sevThreshold = debugLog
	logging.flushInterval = defaultFlushInterval
	logging.enableLogHeader = true
	logging.traceLocation.logger = &logging

}

// StdLogger return internal logging , default logger
func StdLogger() *LoggingT {
	return &logging
}

// NewLogger create new logger
func NewLogger() *LoggingT {
	l := &LoggingT{
		flushInterval: defaultFlushInterval,
		sevThreshold:  debugLog,
	}

	l.traceLocation.logger = l

	return l
}

// ToStderr set logging log to stderr
func (l *LoggingT) ToStderr(enable bool) *LoggingT {
	l.toStderr = enable
	return l
}

// AlsoToStderr set logging log to stderr in addtional
func (l *LoggingT) AlsoToStderr(enable bool) *LoggingT {
	l.alsoToStderr = enable
	return l
}

// SevThreshold set the logging threshold , DEBUG / INFO / WARNING / ERROR / FATAL
func (l *LoggingT) SevThreshold(sevName string) *LoggingT {

	sev := Severity(0)
	err := sev.Set(sevName)
	if err != nil {
		l.exit(err)
	}
	l.sevThreshold.set(sev)
	return l
}

// TraceLocation set the logging traceLocation
func (l *LoggingT) TraceLocation(value string) *LoggingT {
	err := l.traceLocation.Set(value)
	if err != nil {
		l.exit(err)
	}
	return l
}

// LogDir set the logging output dir
func (l *LoggingT) LogDir(dir string) *LoggingT {
	l.logDir = dir
	return l
}

// EnableLogHeader turn on the header for every single line
func (l *LoggingT) EnableLogHeader(enable bool) *LoggingT {
	l.enableLogHeader = enable
	return l
}

// EnableLogLink add symlink to the newest log
func (l *LoggingT) EnableLogLink(enable bool) *LoggingT {
	l.enableLogLink = enable
	return l
}

// HeaderFormat set the logging header format callback
func (l *LoggingT) HeaderFormat(f HeaderFormatFunc) *LoggingT {
	l.headerFormater = f
	return l
}

// StdHeader set the header format as same as LstdFlags in Package `log`
func (l *LoggingT) StdHeader() *LoggingT {

	l.headerFormater = func(buf *bytes.Buffer, _ Severity, ts time.Time, _ int, file string, line int) {
		fmt.Fprintf(buf, "%s %s:%d: ", ts.Format("2006/01/02 15:04:05"), file, line)
	}

	return l
}

// FlushInterval set logger flush interval
func (l *LoggingT) FlushInterval(interval time.Duration) *LoggingT {
	if interval <= 0 {
		interval = defaultFlushInterval
	}
	l.flushInterval = interval
	return l
}

// FileNameFormatFunc tell logger how to format log file name
type FileNameFormatFunc func(severityLevel string, ts time.Time) string

// FileNameFormat set log file name formater
func (l *LoggingT) FileNameFormat(formatFunc FileNameFormatFunc) *LoggingT {
	l.fileNameFormatFunc = formatFunc
	return l
}

// Init mark configuration done
func (l *LoggingT) Init() *Logger {
	l.initOnce.Do(func() {
		go l.flushDaemon()
	})

	return &Logger{l}
}

// Flush flushes all pending log I/O.
func Flush() {
	logging.lockAndFlushAll()
}

func (l *LoggingT) flush() {
	l.lockAndFlushAll()
}

// Flush flushes all pending log I/O.
func (l *Logger) Flush() {
	l.l.lockAndFlushAll()
}

// LoggingT collects all the global state of the logging setup.
type LoggingT struct {
	// Boolean flags. Not handled atomically because the flag.Value interface
	// does not let us avoid the =true, and that shorthand is necessary for
	// compatibility. TODO: does this matter enough to fix? Seems unlikely.
	toStderr     bool // The -logtostderr flag.
	alsoToStderr bool // The -alsologtostderr flag.

	// Level flag. Handled atomically.
	sevThreshold Severity

	// freeList is a list of byte buffers, maintained under freeListMu.
	freeList *buffer
	// freeListMu maintains the free list. It is separate from the main mutex
	// so buffers can be grabbed and printed to without holding the main lock,
	// for better parallelization.
	freeListMu sync.Mutex

	// mu protects the remaining elements of this structure and is
	// used to synchronize logging.
	mu sync.Mutex
	// file holds writer for each of the log types.
	file [numSeverity]flushSyncWriter

	// traceLocation is the state of the -log_backtrace_at flag.
	traceLocation traceLocation

	// to format the custom header
	headerFormater HeaderFormatFunc

	fileNameFormatFunc FileNameFormatFunc

	// If non-empty, overrides the choice of directory in which to write logs.
	// See createLogDirs for the full list of possible destinations.
	// var logDir = flag.String("log_dir", "", "If non-empty, write log files in this directory")
	logDir     string
	onceLogDir sync.Once

	// guard with Init()
	initOnce sync.Once

	// fatalNoStacks is non-zero if we are to exit without dumping goroutine stacks.
	// It allows Exit and relatives to use the Fatal logs.
	fatalNoStacks uint32

	enableLogHeader bool
	enableLogLink   bool // symlink to the newest log with level-tag

	flushInterval time.Duration
}

// Logger is configured and return by Init
type Logger struct {
	l *LoggingT
}

// default logger
var logging LoggingT

// getBuffer returns a new, ready-to-use buffer.
func (l *LoggingT) getBuffer() *buffer {
	l.freeListMu.Lock()
	b := l.freeList
	if b != nil {
		l.freeList = b.next
	}
	l.freeListMu.Unlock()
	if b == nil {
		b = new(buffer)
	} else {
		b.next = nil
		b.Reset()
	}
	return b
}

// putBuffer returns a buffer to the free list.
func (l *LoggingT) putBuffer(b *buffer) {
	if b.Len() >= 256 {
		// Let big buffers die a natural death.
		return
	}
	l.freeListMu.Lock()
	b.next = l.freeList
	l.freeList = b
	l.freeListMu.Unlock()
}

var timeNow = time.Now // Stubbed out for testing.

/*
header formats a log header as defined by the C++ implementation.
It returns a buffer containing the formatted header and the user's file and line number.
The depth specifies how many stack frames above lives the source line to be identified in the log message.

Log lines have this form:
	Lmmdd hh:mm:ss.uuuuuu threadid file:line] msg...
where the fields are defined as follows:
	L                A single character, representing the log level (eg 'I' for INFO)
	mm               The month (zero padded; ie May is '05')
	dd               The day (zero padded)
	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
	threadid         The space-padded thread ID as returned by GetTID()
	file             The file name
	line             The line number
	msg              The user-supplied message
*/
func (l *LoggingT) header(s Severity, depth int) (*buffer, string, int) {
	_, file, line, ok := runtime.Caller(4 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return l.formatHeader(s, file, line), file, line
}

// HeaderFormatFunc define the callback to made up custom header
type HeaderFormatFunc func(buf *bytes.Buffer, l Severity, ts time.Time, pid int, file string, line int)

// formatHeader formats a log header using the provided file name and line number.
func (l *LoggingT) formatHeader(s Severity, file string, line int) *buffer {
	buf := l.getBuffer()
	if !l.enableLogHeader {
		return buf
	}

	now := timeNow()
	if line < 0 {
		line = 0 // not a real line number, but acceptable to someDigits
	}
	if s > fatalLog {
		s = infoLog // for safety.
	}

	// use custom header
	if l.headerFormater != nil {
		l.headerFormater(&buf.Buffer, s, now, pid, file, line)
		return buf
	}

	// Avoid Fprintf, for speed. The format is so simple that we can do it quickly by hand.
	// It's worth about 3X. Fprintf is hard.
	_, month, day := now.Date()
	hour, minute, second := now.Clock()
	// Lmmdd hh:mm:ss.uuuuuu threadid file:line]
	buf.tmp[0] = severityChar[s]
	buf.twoDigits(1, int(month))
	buf.twoDigits(3, day)
	buf.tmp[5] = ' '
	buf.twoDigits(6, hour)
	buf.tmp[8] = ':'
	buf.twoDigits(9, minute)
	buf.tmp[11] = ':'
	buf.twoDigits(12, second)
	buf.tmp[14] = '.'
	buf.nDigits(6, 15, now.Nanosecond()/1000, '0')
	buf.tmp[21] = ' '
	buf.nDigits(7, 22, pid, ' ') // TODO: should be TID
	buf.tmp[29] = ' '
	buf.Write(buf.tmp[:30])
	buf.WriteString(file)
	buf.tmp[0] = ':'
	n := buf.someDigits(1, line)
	buf.tmp[n+1] = ']'
	buf.tmp[n+2] = ' '
	buf.Write(buf.tmp[:n+3])
	return buf
}

func (l *LoggingT) println(s Severity, args ...interface{}) {
	buf, file, line := l.header(s, 0)
	fmt.Fprintln(buf, args...)
	l.output(s, buf, file, line, false)
}

// Print do print depent on given severity
func (l *LoggingT) print(s Severity, args ...interface{}) {
	l.printDepth(s, 1, args...)
}

func (l *LoggingT) printDepth(s Severity, depth int, args ...interface{}) {
	buf, file, line := l.header(s, depth)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, false)
}

// Printf do printf depent on given severity
func (l *LoggingT) printf(s Severity, format string, args ...interface{}) {
	buf, file, line := l.header(s, 0)
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, false)
}

// printWithFileLine behaves like print but uses the provided file and line number.  If
// alsoLogToStderr is true, the log message always appears on standard error; it
// will also appear in the log file unless --logtostderr is set.
func (l *LoggingT) printWithFileLine(s Severity, file string, line int, alsoToStderr bool, args ...interface{}) {
	buf := l.formatHeader(s, file, line)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, alsoToStderr)
}

// output writes the data to the log files and releases the buffer.
func (l *LoggingT) output(s Severity, buf *buffer, file string, line int, alsoToStderr bool) {

	if s < l.sevThreshold.get() {
		return
	}

	l.mu.Lock()
	if l.traceLocation.isSet() {
		if l.traceLocation.match(file, line) {
			buf.Write(stacks(false))
		}
	}
	data := buf.Bytes()
	if l.toStderr {
		os.Stderr.Write(data)
	} else {
		if alsoToStderr || l.alsoToStderr {
			os.Stderr.Write(data)
		}
		if l.file[s] == nil {
			if err := l.createFiles(s); err != nil {
				os.Stderr.Write(data) // Make sure the message appears somewhere.
				l.exit(err)
			}
		}
		switch s {
		case fatalLog:
			l.file[fatalLog].Write(data)
			// fallthrough
		case errorLog:
			l.file[errorLog].Write(data)
			// fallthrough
		case warningLog:
			l.file[warningLog].Write(data)
			// fallthrough
		case infoLog:
			l.file[infoLog].Write(data)
		case debugLog:
			l.file[debugLog].Write(data)
		case accessLog:
			l.file[accessLog].Write(data)
		case interfaceAvgDurationLog:
			l.file[interfaceAvgDurationLog].Write(data)
		}
	}
	if s == fatalLog {
		// If we got here via Exit rather than Fatal, print no stacks.
		if atomic.LoadUint32(&l.fatalNoStacks) > 0 {
			l.mu.Unlock()
			l.timeoutFlush(10 * time.Second)
			os.Exit(1)
		}
		// Dump all goroutine stacks before exiting.
		// First, make sure we see the trace for the current goroutine on standard error.
		// If -logtostderr has been specified, the loop below will do that anyway
		// as the first stack in the full dump.
		if !l.toStderr {
			os.Stderr.Write(stacks(false))
		}
		// Write the stack trace for all goroutines to the files.
		trace := stacks(true)
		logExitFunc = func(error) {} // If we get a write error, we'll still exit below.
		for log := fatalLog; log >= debugLog; log-- {
			if f := l.file[log]; f != nil { // Can be nil if -logtostderr is set.
				f.Write(trace)
			}
		}
		l.mu.Unlock()
		l.timeoutFlush(10 * time.Second)
		os.Exit(255) // C++ uses -1, which is silly because it's anded with 255 anyway.
	}
	l.putBuffer(buf)
	l.mu.Unlock()

}

// timeoutFlush calls Flush and returns when it completes or after timeout
// elapses, whichever happens first.  This is needed because the hooks invoked
// by Flush may deadlock when glog.Fatal is called from a hook that holds
// a lock.
func (l *LoggingT) timeoutFlush(timeout time.Duration) {
	done := make(chan bool, 1)
	go func() {
		l.flush() // calls lockAndFlushAll()
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "glog: Flush took longer than", timeout)
	}
}

// stacks is a wrapper for runtime.Stack that attempts to recover the data for all goroutines.
func stacks(all bool) []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit. Start large, though.
	n := 10000
	if all {
		n = 100000
	}
	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
		n *= 2
	}
	return trace
}

// logExitFunc provides a simple mechanism to override the default behavior
// of exiting on error. Used in testing and to guarantee we reach a required exit
// for fatal logs. Instead, exit could be a function rather than a method but that
// would make its use clumsier.
var logExitFunc func(error) = func(err error) {
	fmt.Fprintf(os.Stderr, "glog: fatal error: %v", err)
}

// exit is called if there is trouble creating or writing log files.
// It flushes the logs and exits the program; there's no point in hanging around.
// l.mu is held.
func (l *LoggingT) exit(err error) {
	fmt.Fprintf(os.Stderr, "log: exiting because of error: %s\n", err)
	// If logExitFunc is set, we do that instead of exiting.
	if logExitFunc != nil {
		logExitFunc(err)
		return
	}
	// l.flushAll()
	// os.Exit(2)
}

// createFiles creates all the log files for severity from sev down to infoLog.
// l.mu is held.
func (l *LoggingT) createFiles(sev Severity) error {
	now := time.Now()
	// Files are created in decreasing severity order, so as soon as we find one
	// has already been created, we can stop.
	for s := sev; s >= debugLog && l.file[s] == nil; s-- {
		sb := &syncBuffer{
			logger: l,
			sev:    s,
		}
		if err := sb.RotateFile(now); err != nil {
			return err
		}
		l.file[s] = sb
	}
	return nil
}

const defaultFlushInterval = 3 * time.Second

// flushDaemon periodically flushes the log file buffers.
func (l *LoggingT) flushDaemon() {
	for _ = range time.NewTicker(l.flushInterval).C {
		l.lockAndFlushAll()
	}
}

// lockAndFlushAll is like flushAll but locks l.mu first.
func (l *LoggingT) lockAndFlushAll() {
	l.mu.Lock()
	l.flushAll()
	l.mu.Unlock()
}

// flushAll flushes all the logs and attempts to "sync" their data to disk.
// l.mu is held.
func (l *LoggingT) flushAll() {
	// Flush from fatal down, in case there's trouble flushing.
	for s := fatalLog; s >= debugLog; s-- {
		file := l.file[s]
		if file != nil {
			file.Flush() // ignore error
			file.Sync()  // ignore error

			if err := file.RotateFile(time.Now()); err != nil {
				l.exit(err)
			}
		}
	}

}

// CopyStandardLogTo arranges for messages written to the Go "log" package's
// default logs to also appear in the Google logs for the named and lower
// severities.  Subsequent changes to the standard log's default output location
// or format may break this behavior.
//
// Valid names are "INFO", "WARNING", "ERROR", and "FATAL".  If the name is not
// recognized, CopyStandardLogTo panics.
func CopyStandardLogTo(name string) {
	sev, ok := severityByName(name)
	if !ok {
		panic(fmt.Sprintf("log.CopyStandardLogTo(%q): unrecognized severity name", name))
	}
	// Set a log format that captures the user's file and line:
	//   d.go:23: message
	stdLog.SetFlags(stdLog.Lshortfile)
	stdLog.SetOutput(logBridge(sev))
}

// logBridge provides the Write method that enables CopyStandardLogTo to connect
// Go's standard logs to the logs provided by this package.
type logBridge Severity

// Write parses the standard logging line and passes its components to the
// logger for severity(lb).
func (lb logBridge) Write(b []byte) (n int, err error) {
	var (
		file = "???"
		line = 1
		text string
	)
	// Split "d.go:23: message" into "d.go", "23", and "message".
	if parts := bytes.SplitN(b, []byte{':'}, 3); len(parts) != 3 || len(parts[0]) < 1 || len(parts[2]) < 1 {
		text = fmt.Sprintf("bad log format: %s", b)
	} else {
		file = string(parts[0])
		text = string(parts[2][1:]) // skip leading space
		line, err = strconv.Atoi(string(parts[1]))
		if err != nil {
			text = fmt.Sprintf("bad line number: %s", b)
			line = 1
		}
	}
	// printWithFileLine with alsoToStderr=true, so standard log messages
	// always appear on standard error.
	logging.printWithFileLine(Severity(lb), file, line, true, text)
	return len(b), nil
}

// Debug logs to the DEBUG log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Debug(args ...interface{}) {
	logging.print(debugLog, args...)
}

// Debug logs to the DEBUG log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func (l *Logger) Debug(args ...interface{}) {
	l.l.print(debugLog, args...)
}

// DebugDepth acts as Debug but uses depth to determine which call frame to log.
// DebugDepth(0, "msg") is the same as Debug("msg").
func DebugDepth(depth int, args ...interface{}) {
	logging.printDepth(debugLog, depth, args...)
}

// DebugDepth acts as Debug but uses depth to determine which call frame to log.
// DebugDepth(0, "msg") is the same as Debug("msg").
func (l *Logger) DebugDepth(depth int, args ...interface{}) {
	l.l.printDepth(debugLog, depth, args...)
}

// Debugln logs to the DEBUG log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Debugln(args ...interface{}) {
	logging.println(debugLog, args...)
}

// Debugln logs to the DEBUG log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func (l *Logger) Debugln(args ...interface{}) {
	l.l.println(debugLog, args...)
}

// Debugf logs to the DEBUG log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Debugf(format string, args ...interface{}) {
	logging.printf(debugLog, format, args...)
}

// Debugf logs to the DEBUG log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.l.printf(debugLog, format, args...)
}

// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Info(args ...interface{}) {
	logging.print(infoLog, args...)
}

// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func (l *Logger) Info(args ...interface{}) {
	l.l.print(infoLog, args...)
}

// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth(0, "msg") is the same as Info("msg").
func InfoDepth(depth int, args ...interface{}) {
	logging.printDepth(infoLog, depth, args...)
}

// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth(0, "msg") is the same as Info("msg").
func (l *Logger) InfoDepth(depth int, args ...interface{}) {
	l.l.printDepth(infoLog, depth, args...)
}

// Infoln logs to the INFO log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Infoln(args ...interface{}) {
	logging.println(infoLog, args...)
}

// Infoln logs to the INFO log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func (l *Logger) Infoln(args ...interface{}) {
	l.l.println(infoLog, args...)
}

// Infof logs to the INFO log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {
	logging.printf(infoLog, format, args...)
}

// Infof logs to the INFO log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.l.printf(infoLog, format, args...)
}

// Warning logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Warning(args ...interface{}) {
	logging.print(warningLog, args...)
}

// Warning logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func (l *Logger) Warning(args ...interface{}) {
	l.l.print(warningLog, args...)
}

// WarningDepth acts as Warning but uses depth to determine which call frame to log.
// WarningDepth(0, "msg") is the same as Warning("msg").
func WarningDepth(depth int, args ...interface{}) {
	logging.printDepth(warningLog, depth, args...)
}

// WarningDepth acts as Warning but uses depth to determine which call frame to log.
// WarningDepth(0, "msg") is the same as Warning("msg").
func (l *Logger) WarningDepth(depth int, args ...interface{}) {
	l.l.printDepth(warningLog, depth, args...)
}

// Warningln logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Warningln(args ...interface{}) {
	logging.println(warningLog, args...)
}

// Warningln logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func (l *Logger) Warningln(args ...interface{}) {
	l.l.println(warningLog, args...)
}

// Warningf logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	logging.printf(warningLog, format, args...)
}

// Warningf logs to the WARNING and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.l.printf(warningLog, format, args...)
}

// Error logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Error(args ...interface{}) {
	logging.print(errorLog, args...)
}

// Error logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func (l *Logger) Error(args ...interface{}) {
	l.l.print(errorLog, args...)
}

// ErrorDepth acts as Error but uses depth to determine which call frame to log.
// ErrorDepth(0, "msg") is the same as Error("msg").
func ErrorDepth(depth int, args ...interface{}) {
	logging.printDepth(errorLog, depth, args...)
}

// ErrorDepth acts as Error but uses depth to determine which call frame to log.
// ErrorDepth(0, "msg") is the same as Error("msg").
func (l *Logger) ErrorDepth(depth int, args ...interface{}) {
	l.l.printDepth(errorLog, depth, args...)
}

// Errorln logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Errorln(args ...interface{}) {
	logging.println(errorLog, args...)
}

// Errorln logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func (l *Logger) Errorln(args ...interface{}) {
	l.l.println(errorLog, args...)
}

// Errorf logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	logging.printf(errorLog, format, args...)
}

// Errorf logs to the ERROR, WARNING, and INFO logs.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.l.printf(errorLog, format, args...)
}

// Fatal logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Fatal(args ...interface{}) {
	logging.print(fatalLog, args...)
}

// Fatal logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func (l *Logger) Fatal(args ...interface{}) {
	l.l.print(fatalLog, args...)
}

// FatalDepth acts as Fatal but uses depth to determine which call frame to log.
// FatalDepth(0, "msg") is the same as Fatal("msg").
func FatalDepth(depth int, args ...interface{}) {
	logging.printDepth(fatalLog, depth, args...)
}

// FatalDepth acts as Fatal but uses depth to determine which call frame to log.
// FatalDepth(0, "msg") is the same as Fatal("msg").
func (l *Logger) FatalDepth(depth int, args ...interface{}) {
	l.l.printDepth(fatalLog, depth, args...)
}

// Fatalln logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Fatalln(args ...interface{}) {
	logging.println(fatalLog, args...)
}

// Fatalln logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func (l *Logger) Fatalln(args ...interface{}) {
	l.l.println(fatalLog, args...)
}

// Fatalf logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Fatalf(format string, args ...interface{}) {
	logging.printf(fatalLog, format, args...)
}

// Fatalf logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.l.printf(fatalLog, format, args...)
}

// Exit logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Exit(args ...interface{}) {
	atomic.StoreUint32(&logging.fatalNoStacks, 1)
	logging.print(fatalLog, args...)
}

// Exit logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func (l *Logger) Exit(args ...interface{}) {
	atomic.StoreUint32(&l.l.fatalNoStacks, 1)
	l.l.print(fatalLog, args...)
}

// ExitDepth acts as Exit but uses depth to determine which call frame to log.
// ExitDepth(0, "msg") is the same as Exit("msg").
func ExitDepth(depth int, args ...interface{}) {
	atomic.StoreUint32(&logging.fatalNoStacks, 1)
	logging.printDepth(fatalLog, depth, args...)
}

// ExitDepth acts as Exit but uses depth to determine which call frame to log.
// ExitDepth(0, "msg") is the same as Exit("msg").
func (l *Logger) ExitDepth(depth int, args ...interface{}) {
	atomic.StoreUint32(&l.l.fatalNoStacks, 1)
	l.l.printDepth(fatalLog, depth, args...)
}

// Exitln logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
func Exitln(args ...interface{}) {
	atomic.StoreUint32(&logging.fatalNoStacks, 1)
	logging.println(fatalLog, args...)
}

// Exitln logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
func (l *Logger) Exitln(args ...interface{}) {
	atomic.StoreUint32(&l.l.fatalNoStacks, 1)
	l.l.println(fatalLog, args...)
}

// Exitf logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Exitf(format string, args ...interface{}) {
	atomic.StoreUint32(&logging.fatalNoStacks, 1)
	logging.printf(fatalLog, format, args...)
}

// Exitf logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func (l *Logger) Exitf(format string, args ...interface{}) {
	atomic.StoreUint32(&l.l.fatalNoStacks, 1)
	l.l.printf(fatalLog, format, args...)
}

//===========
// AccessDepth acts as Acess but uses depth to determine which call frame to log.
// AccessDepth(0, "msg") is the same as Acess("msg").
func (l *Logger) AccessDepth(depth int, args ...interface{}) {
	l.l.printDepth(accessLog, depth, args...)
}

func (l *Logger) InterfaceAvgDurationDepth(depth int, args ...interface{}) {
	l.l.printDepth(interfaceAvgDurationLog, depth, args...)
}
