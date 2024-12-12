package utils

import (
	"fmt"
)

func Info(message string) {
	print(fmt.Sprintf("%s %s@ \n", GetTimeFmtStr(), message))
}

func Error(message string) {
	print(fmt.Sprintf("%s %s× \n", GetTimeFmtStr(), message))
}

func Success(message string) {
	print(fmt.Sprintf("%s %s√ \n", GetTimeFmtStr(), message))
}

func MsgError(msg string) {
	Error(msg)
}

func MsgInfo(msg string) {
	Info(msg)
}

func MsgSuccess(msg string) {
	Success(msg)
}
