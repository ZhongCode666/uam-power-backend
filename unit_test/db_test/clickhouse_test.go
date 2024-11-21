package db

import (
	"fmt"
	"testing"
	dbservice "uam-power-backend/service/db_service"
	DBconfig "uam-power-backend/unit_test/config_test"
)

func TestClickHoseConnect(t *testing.T) {
	cfg := DBconfig.NewConfig()

	ClickHoseService, err := dbservice.NewClickHouse(
		cfg.ClickHouseCfg.Host, cfg.ClickHouseCfg.Port, cfg.ClickHouseCfg.User, cfg.ClickHouseCfg.Psw,
		cfg.ClickHouseCfg.BatchSize, cfg.ClickHouseCfg.FlushPeriod, cfg.ClickHouseCfg.DB, false,
		cfg.ClickHouseCfg.Columns)
	if err != nil {
		t.Error(err)
	}
	err = ClickHoseService.ExecuteCmd("CREATE TABLE IF NOT EXISTS test_go_123 (Longitude Float64 NOT NULL, Latitude Float64 NOT NULL, Altitude Float64 NOT NULL, Yaw Float64 NOT NULL, DataTime DateTime64(6) NOT NULL,  UploadTime DateTime64(6) NOT NULL  DEFAULT now64(6), Event String NOT NULL) ENGINE = MergeTree() ORDER BY DataTime;")
	if err != nil {
		t.Error(err)
	}
	err = ClickHoseService.ExecuteCmd("CREATE TABLE IF NOT EXISTS test22 (Longitude Float64 NOT NULL, Latitude Float64 NOT NULL, Altitude Float64 NOT NULL, Yaw Float64 NOT NULL, DataTime DateTime64(6) NOT NULL,  UploadTime DateTime64(6) NOT NULL DEFAULT now64(6)) ENGINE = MergeTree() ORDER BY DataTime;")
	if err != nil {
		t.Error(err)
	}
	err = ClickHoseService.ExecuteCmd("CREATE TABLE IF NOT EXISTS test11 (DataTime DateTime64(6) NOT NULL,  CreateTime DateTime64(6) NOT NULL DEFAULT now64(6), Event String NOT NULL)  ENGINE = MergeTree() ORDER BY DataTime;")
	if err != nil {
		t.Error(err)
	}

}

func TestClickHoseAddData(t *testing.T) {
	cfg := DBconfig.NewConfig()

	ClickHoseService, err := dbservice.NewClickHouse(
		cfg.ClickHouseCfg.Host, cfg.ClickHouseCfg.Port, cfg.ClickHouseCfg.User, cfg.ClickHouseCfg.Psw,
		cfg.ClickHouseCfg.BatchSize, cfg.ClickHouseCfg.FlushPeriod, cfg.ClickHouseCfg.DB, true,
		cfg.ClickHouseCfg.Columns)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 20; i++ {
		err = ClickHoseService.Add(
			"test_go_123",
			[]string{"Longitude", "Latitude", "Altitude", "Yaw", "DataTime", "Event"},
			[]interface{}{113.2, 22, 11, 150, fmt.Sprintf("2024-11-21 10:11:00.0000%d", i), "test"},
		)
		if err != nil {
			t.Error(err)
		}
	}
}
