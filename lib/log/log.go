package log

import "fmt"

const DEBUG = false

func Info(i ...interface{}) {
	if DEBUG {
		fmt.Println(i...)
	}
}

func Infof(format string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(format+"\n", args...)
	}
}

func Error(i ...interface{}) {
	fmt.Println(i...)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
