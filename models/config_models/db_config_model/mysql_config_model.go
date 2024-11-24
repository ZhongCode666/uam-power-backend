package db_config_model

type MySqlConfigModel struct {
	Usr            string   `yaml:"Usr"`
	Psw            string   `yaml:"Psw"`
	Host           string   `yaml:"Host"`
	DB             string   `yaml:"DB"`
	EventDB        string   `yaml:"EventDB"`
	FlightDB       string   `yaml:"FlightDB"`
	FlightColumn   []string `yaml:"FlightColumn"`
	EventColumn    []string `yaml:"EventColumn"`
	FlightInterval int      `yaml:"FlightInterval"`
	EventInterval  int      `yaml:"EventInterval"`
	Port           int      `yaml:"Port"`
}
