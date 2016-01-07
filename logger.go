package logger

import (
	"github.com/MoonighT/logger/logfile"
	"log"
	"path"
	"runtime"
	"strings"
)

var (
	logger   *log.Logger
	LogLevel int64
)

const (
	DETAIL = "DETAIL"
	INFO   = "INFO"
	WARN   = "WARN"
	ERROR  = "ERROR"
)

func Init(filename string, logLevel, rotateInterval, maxLogFileSize, maxLogFileCount int64) {
	output, err := logfile.Open(filename, rotateInterval, maxLogFileSize, maxLogFileCount)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(output, "", log.LstdFlags|log.Lmicroseconds)
	LogLevel = logLevel
}

func getActualCaller() (file string, line int, ok bool) {
	cpc, _, _, ok := runtime.Caller(2)
	if !ok {
		return
	}

	callerFunPtr := runtime.FuncForPC(cpc)
	if callerFunPtr == nil {
		ok = false
		return
	}

	var pc uintptr
	for callLevel := 3; callLevel < 5; callLevel++ {
		pc, file, line, ok = runtime.Caller(callLevel)
		if !ok {
			return
		}
		funcPtr := runtime.FuncForPC(pc)
		if funcPtr == nil {
			ok = false
			return
		}
		if getFuncNameWithoutPackage(funcPtr.Name()) !=
			getFuncNameWithoutPackage(callerFunPtr.Name()) {
			return
		}
	}
	ok = false
	return
}

func getFuncNameWithoutPackage(name string) string {
	pos := strings.LastIndex(name, ".")
	if pos >= 0 {
		name = name[pos+1:]
	}
	return name
}

func logf(level int, prefix, format string, v ...interface{}) {
	if int(LogLevel) >= level {
		file, line, ok := getActualCaller()
		arg := make([]interface{}, 0)
		arg = append(arg, prefix)
		if ok {
			arg = append(arg, path.Base(file), line)
			arg = append(arg, v...)
			logger.Printf("[%s][%s:%d]"+format, arg...)
		} else {
			arg = append(arg, v...)
			logger.Printf("[%s]"+format, arg...)
		}
	}
}

func Detailf(format string, v ...interface{}) {
	logf(1, DETAIL, format, v...)
}

func Infof(format string, v ...interface{}) {
	logf(0, INFO, format, v...)
}

func Warnf(format string, v ...interface{}) {
	logf(0, WARN, format, v...)
}

func Errorf(format string, v ...interface{}) {
	logf(0, ERROR, format, v...)
}
