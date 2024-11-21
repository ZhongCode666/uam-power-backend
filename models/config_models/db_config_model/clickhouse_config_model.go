package db_config_model

type ClickHouseConfigModel struct {
	Host             string   `yaml:"Host"`
	Port             int      `yaml:"Port"`
	Username         string   `yaml:"Username"`
	Password         string   `yaml:"Password"`
	FlightDatabase   string   `yaml:"FlightDatabase"`
	EventDatabase    string   `yaml:"EventDatabase"`
	BatchSize        int      `yaml:"BatchSize"`
	FlushPeriod      int      `yaml:"FlushPeriod"`
	EventBatchSize   int      `yaml:"EventBatchSize"`
	EventFlushPeriod int      `yaml:"EventFlushPeriod"`
	FlightColumn     []string `yaml:"FlightColumn"`
	EventColumn      []string `yaml:"EventColumn"`
}
