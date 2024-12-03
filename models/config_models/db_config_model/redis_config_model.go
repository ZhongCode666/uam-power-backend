package db_config_model

type RedisConfigModel struct {
	Port         int    `yaml:"Port"`         // 端口号
	StatusDBno   int    `yaml:"StatusDBno"`   // 状态数据库编号
	EventDBno    int    `yaml:"EventDBno"`    // 事件数据库编号
	Host         string `yaml:"Host"`         // 主机地址
	AircraftDBno int    `yaml:"AircraftDBno"` // 飞机数据库编号
	TaskInfoDBno int    `yaml:"TaskInfoDBno"` // 任务信息数据库编号
}
