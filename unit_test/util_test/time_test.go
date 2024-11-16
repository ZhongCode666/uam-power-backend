package util

import (
	"testing"
	"uam-power-backend/utils"
)

func TestStrToMysqlTimeString(t *testing.T) {
	t.Log(utils.IsValidSqlTimeFormat("20060102150405000000"))
	t.Log(utils.IsValidSqlTimeFormat("2024-11-16 12:17:00.123456"))
}
