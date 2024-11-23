package db_config_model

type MongoConfigModel struct {
	Usr    string `yaml:"Usr"`
	Psw    string `yaml:"Psw"`
	Host   string `yaml:"Host"`
	AreaDB string `yaml:"AreaDB"`
	Port   int    `yaml:"Port"`
}
