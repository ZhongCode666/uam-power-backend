package db_config_model

// MongoConfigModel 结构体表示 MongoDB 的配置模型
type MongoConfigModel struct {
	Usr    string `yaml:"Usr"`    // 用户名
	Psw    string `yaml:"Psw"`    // 密码
	Host   string `yaml:"Host"`   // 主机地址
	AreaDB string `yaml:"AreaDB"` // 区域数据库
	Port   int    `yaml:"Port"`   // 端口号
}
