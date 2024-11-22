package db_config_model

type MongoConfigModel struct {
	Usr        string `yaml:"Usr"`
	Psw        string `yaml:"Psw"`
	Host       string `yaml:"Host"`
	LaneDataDB string `yaml:"LaneDataDB"`
	AreaDB     string `yaml:"AreaDB"`
	RasterDB   string `yaml:"RasterDB"`
	Port       int    `yaml:"Port"`
}
