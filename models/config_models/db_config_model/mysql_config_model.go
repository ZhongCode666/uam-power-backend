package db_config_model

type MySqlConfigModel struct {
	Usr            string   `yaml:"Usr"`            // 用户名
	Psw            string   `yaml:"Psw"`            // 密码
	Host           string   `yaml:"Host"`           // 主机地址
	DB             string   `yaml:"DB"`             // 数据库
	EventDB        string   `yaml:"EventDB"`        // 事件数据库
	FlightDB       string   `yaml:"FlightDB"`       // 航班数据库
	FlightColumn   []string `yaml:"FlightColumn"`   // 航班列
	EventColumn    []string `yaml:"EventColumn"`    // 事件列
	FlightInterval int      `yaml:"FlightInterval"` // 航班间隔
	EventInterval  int      `yaml:"EventInterval"`  // 事件间隔
	Port           int      `yaml:"Port"`           // 端口号
}
