package DBconfig

type Config struct {
	MySQLCfg struct {
		Usr  string
		Psw  string
		Host string
		Port int
		DB   string
	}
	RedisCfg struct {
		Port int
		DBno int
		Host string
	}
	KafkaAddr  string
	KafkaTopic string
	MongoCfg   struct {
		Host string
		Port int
		DB   string
	}
}

func NewConfig() *Config {
	cfg := Config{
		KafkaAddr:  "175.178.125.164:9092",
		KafkaTopic: "goTest",
	}
	cfg.RedisCfg.Host = "119.29.181.98"
	cfg.RedisCfg.Port = 6379
	cfg.RedisCfg.DBno = 2

	cfg.MySQLCfg.Usr = "keane"
	cfg.MySQLCfg.Psw = "Zhonglaoshizhen6!"
	cfg.MySQLCfg.Host = "gz-cdb-1n3w64q3.sql.tencentcdb.com"
	cfg.MySQLCfg.DB = "node_db"
	cfg.MySQLCfg.Port = 25100

	cfg.MongoCfg.Host = "175.178.125.164"
	cfg.MongoCfg.Port = 27017
	cfg.MongoCfg.DB = "test_go"

	return &cfg
}
