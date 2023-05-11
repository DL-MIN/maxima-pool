/*******************************************************************************
 * Logger
 *
 * @author     Lars Thoms <lars@thoms.io>
 * @date       2022-12-21
 ******************************************************************************/

package main

import (
	"fmt"
	"log"
	"os"
)

const (
	Debug = 0 + iota
	Info
	Warn
	Fatal
)

type Logger struct {
	level      int
	loggerList [4]*log.Logger
}

func (l *Logger) Level() int {
	return l.level
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func NewLogger() *Logger {
	return &Logger{
		level: Debug,
		loggerList: [4]*log.Logger{
			log.New(os.Stdout, "[\u001B[0;37mDEBUG\u001B[0m] ", log.Ldate|log.Ltime|log.Lshortfile),
			log.New(os.Stdout, "[\u001B[0;32mINFO\u001B[0m] ", log.Ldate|log.Ltime|log.Lshortfile),
			log.New(os.Stderr, "[\u001B[0;33mWARNING\u001B[0m] ", log.Ldate|log.Ltime|log.Lshortfile),
			log.New(os.Stderr, "[\u001B[0;31mFATAL\u001B[0m] ", log.Ldate|log.Ltime|log.Lshortfile),
		},
	}
}

func (l *Logger) logWithLevel(level int, format *string, v ...any) {
	if l.level > level {
		return
	}

	var err error

	if format == nil {
		err = l.loggerList[level].Output(3, fmt.Sprint(v...))
	} else {
		err = l.loggerList[level].Output(3, fmt.Sprintf(*format, v...))
	}

	if err != nil {
		fmt.Println(err)
	}

	if level == Fatal {
		os.Exit(1)
	}
}

func (l *Logger) Debug(v ...any) {
	l.logWithLevel(Debug, nil, v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.logWithLevel(Debug, &format, v...)
}

func (l *Logger) Info(v ...any) {
	l.logWithLevel(Info, nil, v...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.logWithLevel(Info, &format, v...)
}

func (l *Logger) Warn(v ...any) {
	l.logWithLevel(Warn, nil, v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.logWithLevel(Warn, &format, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.logWithLevel(Fatal, nil, v...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.logWithLevel(Fatal, &format, v...)
}
