package log

import (
	"fmt"
	"gogotrainschedule/lib/env"
)

const pre_d = "DEBUG: "
const pre_i = "INFO: "
const pre_e = "ERROR: "

func Debug(i ...interface{}) {
	if env.NotProd() {
		i = append([]interface{}{pre_d}, i...)
		fmt.Println(i...)
	}
}

func Debugf(format string, args ...interface{}) {
	if env.NotProd() {
		fmt.Printf(pre_d+format+"\n", args...)
	}
}

func Info(i ...interface{}) {
	i = append([]interface{}{pre_i}, i...)
	fmt.Println(i...)
}

func Infof(format string, args ...interface{}) {
	fmt.Printf(pre_i+format+"\n", args...)
}

func Error(i ...interface{}) {
	i = append([]interface{}{pre_e}, i...)
	fmt.Println(i...)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(pre_e+format+"\n", args...)
}
