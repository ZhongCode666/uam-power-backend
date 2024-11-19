package utils

import (
	"time"
)

func GetTimeStr() string {
	currentTime := time.Now()
	// 定义格式化样式
	formattedTime := currentTime.Format("20060102150405012345")
	return formattedTime
}

func GetTimeFmtStr() string {
	currentTime := time.Now()
	// 定义格式化样式
	formattedTime := currentTime.Format("[2006-01-02 15:04:05.012345]")
	return formattedTime
}

func GetMySqlTimeStr() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05.000000")
}

func TransferTimeStrToSqlTimeStr(str string) string {
	parsedTime, _ := time.Parse("2006-01-02 15:04:05.000000", str)
	// 格式化为目标格式
	return parsedTime.Format("2006-01-02 15:04:05.000000")
}

func IsValidSqlTimeFormat(str string) bool {
	// Parse the string using the desired format
	_, err := time.Parse("2006-01-02 15:04:05.000000", str)
	return err == nil // If err is nil, the string matches the format
}
