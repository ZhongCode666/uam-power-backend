package db_test

import (
	"testing"
	//"github.com/gin-gonic/gin"
	"uam-power-backend/service/db_service"
	"uam-power-backend/unit_test/config_test"
)

func TestRedisConnect(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	_ = dbservice.NewRedisDict(cfg.RedisCfg.Host, cfg.RedisCfg.Port, cfg.RedisCfg.DBno)
}

func TestRedisSet(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	redisFun := dbservice.NewRedisDict(cfg.RedisCfg.Host, cfg.RedisCfg.Port, cfg.RedisCfg.DBno)
	err := redisFun.Set("test_go_str", "test_string")
	if err != nil {
		t.Errorf("test_go_str%s", err)
	}
	err = redisFun.Set("test_go_int", 5)
	if err != nil {
		t.Errorf("test_go_int%s", err)
	}
	err = redisFun.Set("test_go_float", 8.98)
	if err != nil {
		t.Errorf("test_go_float%s", err)
	}
	err = redisFun.Set("test_go_bool", false)
	if err != nil {
		t.Errorf("test_go_bool%s", err)
	}
	var arr1 [5]int
	err = redisFun.Set("test_go_int_list", arr1)
	arr2 := [5]float32{2.2, 3.3, 3.4, 9.9, 8.8}
	if err != nil {
		t.Errorf("test_go_int_list%s", err)
	}
	err = redisFun.Set("test_go_float_list", arr2)
	if err != nil {
		t.Errorf("test_go_float_list%s", err)
	}
	arr3 := [5]string{"SYSU", "test", "go", "zrx", "znb"}
	err = redisFun.Set("test_go_str_list", arr3)
	if err != nil {
		t.Errorf("test_go_str_list%s", err)
	}
	jsons := map[string]interface{}{
		"id":        1,
		"name":      "John Doe",
		"is_active": true,
		"age":       30,
	}
	err = redisFun.Set("test_go_json", jsons)
	if err != nil {
		t.Errorf("test_go_json%s", err)
	}
}

func TestRedisGet(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	redisFun := dbservice.NewRedisDict(cfg.RedisCfg.Host, cfg.RedisCfg.Port, cfg.RedisCfg.DBno)
	re, err := redisFun.Get("test_go_str")
	if err != nil {
		t.Errorf("test_go_str%s", err)
	}
	t.Log(re)
	re, err = redisFun.Get("test_go_int")
	if err != nil {
		t.Errorf("test_go_int%s", err)
	}
	t.Log(re)
	re, err = redisFun.Get("test_go_float")
	if err != nil {
		t.Errorf("test_go_float%s", err)
	}
	t.Log(re)
	re, err = redisFun.Get("test_go_bool")
	if err != nil {
		t.Errorf("test_go_bool%s", err)
	}
	t.Log(re)
	re, err = redisFun.Get("test_go_int_list")
	if err != nil {
		t.Errorf("test_go_int_list%s", err)
	}
	t.Log(re)
	re, err = redisFun.Get("test_go_float_list")
	if err != nil {
		t.Errorf("test_go_float_list%s", err)
	}
	t.Log(re)
	re, err = redisFun.Get("test_go_str_list")
	if err != nil {
		t.Errorf("test_go_str_list%s", err)
	}
	t.Log(re)
	re, err = redisFun.Get("test_go_json")
	if err != nil {
		t.Errorf("test_go_json%s", err)
	}
	t.Log(re)
}

func TestRedisDel(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	redisFun := dbservice.NewRedisDict(cfg.RedisCfg.Host, cfg.RedisCfg.Port, cfg.RedisCfg.DBno)
	err := redisFun.Delete("test_go_str")
	if err != nil {
		t.Errorf("test_delete%s", err)
	}
}

func TestRedisExist(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	redisFun := dbservice.NewRedisDict(cfg.RedisCfg.Host, cfg.RedisCfg.Port, cfg.RedisCfg.DBno)
	re, err := redisFun.Exists("test_go_str")
	if err != nil {
		t.Errorf("test_delete%s", err)
	}
	t.Log(re)
	re, err = redisFun.Exists("test_go_int")
	if err != nil {
		t.Errorf("test_delete%s", err)
	}
	t.Log(re)
}

func TestRedisKeys(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	redisFun := dbservice.NewRedisDict(cfg.RedisCfg.Host, cfg.RedisCfg.Port, cfg.RedisCfg.DBno)
	re, err := redisFun.Keys()
	if err != nil {
		t.Errorf("test_delete%s", err)
	}
	t.Log(re)
}

func TestRedisGetVals(t *testing.T) {
	cfg := DBconfig.NewConfig()
	redisFun := dbservice.NewRedisDict(cfg.RedisCfg.Host, cfg.RedisCfg.Port, cfg.RedisCfg.DBno)
	re, err := redisFun.GetVals([]string{"a", "b", "c"})
	if err != nil {
		t.Errorf("test_getmany%s", err)
	}
	t.Log(re)
}
