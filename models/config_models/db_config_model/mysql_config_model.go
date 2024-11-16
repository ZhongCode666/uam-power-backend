package db_config_model

type MySqlConfigModel struct {
	Usr      string `yaml:"Usr"`
	Psw      string `yaml:"Psw"`
	Host     string `yaml:"Host"`
	DB       string `yaml:"DB"`
	EventDB  string `yaml:"EventDB"`
	FlightDB string `yaml:"FlightDB"`
	Port     int    `yaml:"Port"`
}
