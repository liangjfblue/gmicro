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

// File I/O for logs.

package glog

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func (l *LoggingT) createLogDirs() {
	if l.logDir != "" {
		err := os.MkdirAll(filepath.Clean(l.logDir), os.ModePerm|os.ModeDir)
		if err != nil {
			l.exit(err)
		}
	}

}

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
)

func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

// logName returns a new log file name containing tag, with start time t, and
// the name for the symlink for tag.
func (l *LoggingT) logName(tag string, t time.Time) (name, link string) {

	if l.fileNameFormatFunc != nil {
		name = l.fileNameFormatFunc(tag, t)
	} else {
		name = fmt.Sprintf("%s.%s.%s.log.%s.%04d%02d%02d",
			program,
			host,
			userName,
			tag,
			t.Year(),
			t.Month(),
			t.Day())
	}

	return name, program + "." + tag
}

// create creates a new log file and returns the file and its filename, which
// contains tag ("INFO", "FATAL", etc.) and t.  If the file is created
// successfully, create also attempts to update the symlink for that tag, ignoring
// errors.
func (l *LoggingT) create(tag string, t time.Time) (f *os.File, filename string, err error) {
	l.onceLogDir.Do(func() {
		l.createLogDirs()
	})

	name, link := l.logName(tag, t)
	fname := filepath.Join(l.logDir, name)
	f, err = os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err == nil {
		if l.enableLogLink {
			symlink := filepath.Join(l.logDir, link)
			os.Remove(symlink)        // ignore err
			os.Symlink(name, symlink) // ignore err
		}
		return f, fname, nil
	}

	return nil, "", fmt.Errorf("log: cannot create log: %v", err)

}
