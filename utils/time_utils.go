package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// GetUniqueStr 生成一个唯一的字符串
// 返回一个唯一的十六进制字符串
func GetUniqueStr() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	// 转换为十六进制字符串
	uniqueStr := hex.EncodeToString(bytes)
	return uniqueStr
}

// GetTimeStr 获取当前时间并格式化为特定格式的字符串
// 返回一个格式化后的时间字符串
func GetTimeStr() string {
	currentTime := time.Now()
	// 定义格式化样式
	formattedTime := currentTime.Format("20060102150405012345")
	return formattedTime
}

// GetTimeFmtStr 获取当前时间并格式化为特定格式的字符串
// 返回一个格式化后的时间字符串
func GetTimeFmtStr() string {
	currentTime := time.Now()
	// 定义格式化样式
	formattedTime := currentTime.Format("[2006-01-02 15:04:05.012345]")
	return formattedTime
}

// GetMySqlTimeStr 获取当前时间并格式化为 MySQL 时间格式字符串
// 返回一个符合 MySQL 时间格式的字符串
func GetMySqlTimeStr() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05.000000")
}

// TransferTimeStrToSqlTimeStr 将时间字符串转换为 SQL 时间格式字符串
// str 是要转换的时间字符串
// 返回一个符合 SQL 时间格式的字符串
func TransferTimeStrToSqlTimeStr(str string) string {
	parsedTime, _ := time.Parse("2006-01-02 15:04:05.000000", str)
	// 格式化为目标格式
	return parsedTime.Format("2006-01-02 15:04:05.000000")
}

// IsValidSqlTimeFormat 检查字符串是否符合 SQL 时间格式
// str 是要检查的字符串
// 返回一个布尔值，表示字符串是否符合 SQL 时间格式
func IsValidSqlTimeFormat(str string) bool {
	// 使用指定格式解析字符串
	_, err := time.Parse("2006-01-02 15:04:05.000000", str)
	return err == nil // 如果 err 为 nil，表示字符串符合格式
}
