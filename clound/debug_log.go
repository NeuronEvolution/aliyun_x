package clound

import "fmt"

var debugEnabled = false

func SetDebug(v bool) {
	debugEnabled = v
}

func debugLog(format string, a ...interface{}) {
	if debugEnabled {
		fmt.Printf(format+"\n", a...)
	}
}

func debugWrite(format string, a ...interface{}) {
	if debugEnabled {
		fmt.Printf(format, a...)
	}
}
