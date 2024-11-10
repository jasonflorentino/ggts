package log

import "fmt"

func Info(i ...interface{}) {
	fmt.Println(i...)
}

func Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Error(i ...interface{}) {
	fmt.Println(i...)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
