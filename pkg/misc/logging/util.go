package logging

import (
	"unsafe"
)

func Debug(log DebugLogger, args ...interface{}) {
	if !isNilValue(log) {
		log.Debug(args...)
	}
}

func Debugf(log DebugLogger, template string, args ...interface{}) {
	if !isNilValue(log) {
		log.Debugf(template, args...)
	}
}

func Debugw(log DebugLogger, msg string, keysAndValues ...interface{}) {
	if !isNilValue(log) {
		log.Debugw(msg, keysAndValues...)
	}
}

func Info(log InfoLogger, args ...interface{}) {
	if !isNilValue(log) {
		log.Info(args...)
	}
}

func Infof(log InfoLogger, template string, args ...interface{}) {
	if !isNilValue(log) {
		log.Infof(template, args...)
	}
}

func Infow(log InfoLogger, msg string, keysAndValues ...interface{}) {
	if !isNilValue(log) {
		log.Infow(msg, keysAndValues...)
	}
}

func Error(log ErrorLogger, args ...interface{}) {
	if !isNilValue(log) {
		log.Error(args...)
	}
}

func Errorf(log ErrorLogger, template string, args ...interface{}) {
	if !isNilValue(log) {
		log.Errorf(template, args...)
	}
}

func Errorw(log ErrorLogger, msg string, keysAndValues ...interface{}) {
	if !isNilValue(log) {
		log.Errorw(msg, keysAndValues...)
	}
}

func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}
