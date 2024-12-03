package db_config_model

// ClickHouseConfigModel 结构体表示 ClickHouse 数据库的配置模型
type ClickHouseConfigModel struct {
	Host             string   `yaml:"Host"`             // 主机地址
	Port             int      `yaml:"Port"`             // 端口号
	Username         string   `yaml:"Username"`         // 用户名
	Password         string   `yaml:"Password"`         // 密码
	FlightDatabase   string   `yaml:"FlightDatabase"`   // 航班数据库
	EventDatabase    string   `yaml:"EventDatabase"`    // 事件数据库
	BatchSize        int      `yaml:"BatchSize"`        // 批处理大小
	FlushPeriod      int      `yaml:"FlushPeriod"`      // 刷新周期
	EventBatchSize   int      `yaml:"EventBatchSize"`   // 事件批处理大小
	EventFlushPeriod int      `yaml:"EventFlushPeriod"` // 事件刷新周期
	FlightColumn     []string `yaml:"FlightColumn"`     // 航班列
	EventColumn      []string `yaml:"EventColumn"`      // 事件列
}
