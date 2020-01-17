package glog

import (
	"strconv"
	"strings"
	"sync/atomic"
)

// Severity identifies the sort of log: info, warning etc. It also implements
// the flag.Value interface. The -stderrthreshold flag is of type Severity and
// should be modified only through the flag.Value interface. The values match
// the corresponding constants in C++.
type Severity int32 // sync/atomic int32

// These constants identify the log levels in order of increasing severity.
// A message written to a high-severity log file is also written to each
// lower-severity log file.
const (
	debugLog Severity = iota
	infoLog
	warningLog
	errorLog
	accessLog
	interfaceAvgDurationLog
	fatalLog
	numSeverity = 7
)

// These constants identify the log levels in order of increasing severity.
const (
	// AllLog   = debugLog - 1
	DebugLog                = debugLog
	InfoLog                 = infoLog
	WarnLog                 = warningLog
	ErrorLog                = errorLog
	AccessLog               = accessLog
	InterfaceAvgDurationLog = interfaceAvgDurationLog
	FatalLog                = fatalLog
)

const severityChar = "DIWEF"

// Severity name , options for Sevthreshold
const (
	SevDebug                = "DEBUG"
	SevInfo                 = "INFO"
	SevWarn                 = "WARNING"
	SevError                = "ERROR"
	SevAccess               = "ACCESS"
	SevInterfaceAvgDuration = "IAVGD"
	SevFatal                = "FATAL"
)

var severityName = []string{
	debugLog:                SevDebug,
	infoLog:                 SevInfo,
	warningLog:              SevWarn,
	errorLog:                SevError,
	accessLog:               SevAccess,
	interfaceAvgDurationLog: SevInterfaceAvgDuration,
	fatalLog:                SevFatal,
}

// get returns the value of the severity.
func (s *Severity) get() Severity {
	return Severity(atomic.LoadInt32((*int32)(s)))
}

// set sets the value of the severity.
func (s *Severity) set(val Severity) {
	atomic.StoreInt32((*int32)(s), int32(val))
}

// String is part of the flag.Value interface.
func (s *Severity) String() string {
	return strconv.FormatInt(int64(*s), 10)
}

// Get is part of the flag.Value interface.
func (s *Severity) Get() interface{} {
	return *s
}

// Set is part of the flag.Value interface.
func (s *Severity) Set(value string) error {

	// Is it a known name?
	if v, ok := severityByName(value); ok {
		*s = v
	} else {
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*s = Severity(v)
	}

	return nil
}

func severityByName(s string) (Severity, bool) {
	s = strings.ToUpper(s)
	for i, name := range severityName {
		if name == s {
			return Severity(i), true
		}
	}
	return 0, false
}
