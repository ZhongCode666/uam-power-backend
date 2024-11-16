package db_test

import (
	"fmt"
	"testing"
	//"github.com/gin-gonic/gin"
	"uam-power-backend/service/db_service"
	"uam-power-backend/unit_test/config_test"
)

func TestMySqlAddTable(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLCfg.Usr, cfg.MySQLCfg.Psw, cfg.MySQLCfg.Host, cfg.MySQLCfg.Port,
		cfg.MySQLCfg.DB,
	)
	mysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		panic("Failed to initialize MySQL: " + err.Error())
	}
	_, err = mysqlService.ExecuteCmd("CREATE TABLE go_test (Upload_time DATETIME(6) NOT NULL PRIMARY KEY,Locate_time  DATETIME(6),gga_mode CHAR(10), longitude DOUBLE(15, 12),latitude DOUBLE(15, 12),star_num TINYINT UNSIGNED);")
	if err != nil {
		return
	}
}

func TestMySqlDropTable(t *testing.T) {
	cfg := DBconfig.NewConfig()
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLCfg.Usr, cfg.MySQLCfg.Psw, cfg.MySQLCfg.Host, cfg.MySQLCfg.Port,
		cfg.MySQLCfg.DB,
	)
	mysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		panic("Failed to initialize MySQL: " + err.Error())
	}
	_, err = mysqlService.ExecuteCmd("Drop TABLE go_test;")
	if err != nil {
		return
	}
}

func TestMySqlAddData(t *testing.T) {
	cfg := DBconfig.NewConfig()
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLCfg.Usr, cfg.MySQLCfg.Psw, cfg.MySQLCfg.Host, cfg.MySQLCfg.Port,
		cfg.MySQLCfg.DB,
	)
	mysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		panic("Failed to initialize MySQL: " + err.Error())
	}
	_, err = mysqlService.ExecuteCmd("INSERT INTO go_test (upload_time, locate_time, gga_mode, longitude, latitude, star_num) VALUES ('2024-11-11 11:13:00.000000', '2024-11-11 11:13:00.000000', 'GPS', 113, 22, 22)")
	if err != nil {
		t.Errorf(`%s`, err)
		return
	}
}

func TestMySqlQueryRow(t *testing.T) {
	cfg := DBconfig.NewConfig()
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLCfg.Usr, cfg.MySQLCfg.Psw, cfg.MySQLCfg.Host, cfg.MySQLCfg.Port,
		cfg.MySQLCfg.DB,
	)
	mysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		panic("Failed to initialize MySQL: " + err.Error())
	}
	re, err := mysqlService.QueryRow("Select * from go_test;")
	if err != nil {
		t.Errorf(`%s`, err)
	}
	t.Log(re)
}

func TestMySqlQueryRows(t *testing.T) {
	cfg := DBconfig.NewConfig()
	mysqlLink := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLCfg.Usr, cfg.MySQLCfg.Psw, cfg.MySQLCfg.Host, cfg.MySQLCfg.Port,
		cfg.MySQLCfg.DB,
	)
	mysqlService, err := dbservice.NewMySQLService(mysqlLink)
	if err != nil {
		panic("Failed to initialize MySQL: " + err.Error())
	}
	re, err := mysqlService.QueryRows("Select * from go_test;")
	if err != nil {
		t.Errorf(`%s`, err)
	}
	t.Log(re)
}
