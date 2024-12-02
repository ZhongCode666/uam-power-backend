package utils

import (
	"gopkg.in/yaml.v3"
	"os"
	"uam-power-backend/models/config_models/db_config_model"
	"uam-power-backend/models/config_models/global_model_config_model"
)

// LoadDBConfig 加载数据库配置
// filename 是配置文件的路径
// 返回一个 DbConfigModel 指针和一个错误
func LoadDBConfig(filename string) (*db_config_model.DbConfigModel, error) {
	data, err := os.ReadFile(filename) // 读取文件内容
	if err != nil {
		return nil, err // 如果读取文件出错，返回错误
	}

	var config db_config_model.DbConfigModel
	err = yaml.Unmarshal(data, &config) // 解析 YAML 数据
	if err != nil {
		return nil, err // 如果解析出错，返回错误
	}

	return &config, nil // 返回解析后的配置
}

// LoadGlobalConfig 加载全局配置
// filename 是配置文件的路径
// 返回一个 GlobalConfig 指针和一个错误
func LoadGlobalConfig(filename string) (*global_model_config_model.GlobalConfig, error) {
	data, err := os.ReadFile(filename) // 读取文件内容
	if err != nil {
		return nil, err // 如果读取文件出错，返回错误
	}

	var config global_model_config_model.GlobalConfig
	err = yaml.Unmarshal(data, &config) // 解析 YAML 数据
	if err != nil {
		return nil, err // 如果解析出错，返回错误
	}

	return &config, nil // 返回解析后的配置
}
