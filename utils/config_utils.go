package utils

import (
	"gopkg.in/yaml.v3"
	"os"
	"uam-power-backend/models/config_models/db_config_model"
)

func LoadDBConfig(filename string) (*db_config_model.DbConfigModel, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config db_config_model.DbConfigModel
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
