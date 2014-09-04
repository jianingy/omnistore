/*
 * filename   : log.go
 * created at : 2014-07-20 09:45:13
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package utils

import (
    "bytes"
    "fmt"
    "os"
)

type severity int32 // sync/atomic int32

const (
	debugLog severity = iota
    infoLog
	warningLog
	errorLog
	fatalLog
)

var severityName = []string{
	debugLog:   "DEBUG",
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

func print(s severity, args ...interface{}) {
    var buf bytes.Buffer
	fmt.Fprintf(&buf, "[%s] ", severityName[s])
    if len(args) > 1 {
        fmt.Fprintf(&buf, args[0].(string), args[1:]...)
    } else {
        fmt.Fprint(&buf, args...)
    }
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	fmt.Fprint(os.Stdout, buf.String())
}

func Info(args ...interface{}) {
    print(infoLog, args...)
}

func Debug(args ...interface{}) {
    print(debugLog, args...)
}

func Warn(args ...interface{}) {
    print(warningLog, args...)
}

func Error(args ...interface{}) {
    print(errorLog, args...)
}

func Fatal(args ...interface{}) {
    print(fatalLog, args...)
    panic("fatal error")
}
