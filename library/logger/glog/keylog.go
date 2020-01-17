package glog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// KeyLogger 根据key将记录写入不同文件 , key即文件名
type KeyLogger struct {
	keysMu sync.Mutex
	keys   map[string]*keyWriter

	logdir string

	flushInterval time.Duration
	timeoutClose  time.Duration
}

func NewKeyLogger(logdir string, flushInterval time.Duration, timeoutClose time.Duration) *KeyLogger {
	l := &KeyLogger{
		keys:          make(map[string]*keyWriter),
		logdir:        logdir,
		flushInterval: flushInterval,
		timeoutClose:  timeoutClose,
	}

	go l.flushTimer()

	return l
}

func (l *KeyLogger) flushTimer() {

	for _ = range time.NewTicker(l.flushInterval).C {

		timeout := time.Now().Unix() - int64(l.timeoutClose/1000000000)

		l.keysMu.Lock()
		for _, w := range l.keys {

			if w.touchTime <= timeout {
				w.close()
			} else {
				w.writer.Flush() // ignore error
			}

		}
		l.keysMu.Unlock()
	}

}

func (l *KeyLogger) Flush() {
	l.keysMu.Lock()
	for _, w := range l.keys {
		w.writer.Flush() // ignore error
	}
	l.keysMu.Unlock()
}

func (l *KeyLogger) KPrintf(key string, format string, args ...interface{}) error {

	l.keysMu.Lock()
	defer l.keysMu.Unlock()

	if _, ok := l.keys[key]; !ok {
		l.keys[key] = &keyWriter{k: key, logger: l}

		if err := l.keys[key].createFile(); err != nil {
			return err
		}
	}

	l.keys[key].touchTime = time.Now().Unix()
	fmt.Fprintf(l.keys[key].writer, format, args...)

	return nil

}

func (l *KeyLogger) KPrint(key string, args ...interface{}) error {

	l.keysMu.Lock()
	defer l.keysMu.Unlock()

	if _, ok := l.keys[key]; !ok {
		l.keys[key] = &keyWriter{k: key, logger: l}

		if err := l.keys[key].createFile(); err != nil {
			return err
		}
	}

	l.keys[key].touchTime = time.Now().Unix()
	fmt.Fprint(l.keys[key].writer, args...)

	return nil

}

type keyWriter struct {
	logger *KeyLogger
	k      string

	// bufMu     sync.Mutex
	writer    *bufio.Writer
	f         *os.File
	touchTime int64
}

func (b *keyWriter) createFile() error {

	if err := os.MkdirAll(filepath.Clean(b.logger.logdir), os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(b.logger.logdir, string(b.k)), os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "keylog: createFile: %v", err) // in case we ignore this error
		return err
	}

	b.f = f
	b.writer = bufio.NewWriter(f)

	return nil
}

func (b *keyWriter) close() error {
	err := b.writer.Flush()
	if err != nil {
		return err
	}

	err = b.f.Sync()
	if err != nil {
		return err
	}

	// remove for parent logger
	delete(b.logger.keys, b.k)

	return b.f.Close()
}
