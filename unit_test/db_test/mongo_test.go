package db_test

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	//"github.com/gin-gonic/gin"
	"uam-power-backend/service/db_service"
	"uam-power-backend/unit_test/config_test"
)

func TestMongoAddCollection(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%d",
		cfg.MongoCfg.Host, cfg.MongoCfg.Port,
	)
	mongoService, err := dbservice.NewMongoDBClient(mongoLink, cfg.MongoCfg.DB)
	if err != nil {
		panic("Failed to initialize Mongo: " + err.Error())
	}
	err = mongoService.CreateCollection("go_test")
	if err != nil {
		return
	}
}

func TestMongoAddData(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%d",
		cfg.MongoCfg.Host, cfg.MongoCfg.Port,
	)
	mongoService, err := dbservice.NewMongoDBClient(mongoLink, cfg.MongoCfg.DB)
	if err != nil {
		panic("Failed to initialize Mongo: " + err.Error())
	}
	document := bson.M{"name": "test2", "age": 35}
	re, err := mongoService.InsertOne("go_test", document)
	if err != nil {
		return
	}
	t.Log(re)
}

func TestMongoFindOneData(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%d",
		cfg.MongoCfg.Host, cfg.MongoCfg.Port,
	)
	mongoService, err := dbservice.NewMongoDBClient(mongoLink, cfg.MongoCfg.DB)
	if err != nil {
		panic("Failed to initialize Mongo: " + err.Error())
	}
	filter := bson.M{
		"age": bson.M{"$gt": 29}, // $gt 代表大于
	}
	re, err := mongoService.FindOne("go_test", filter)
	if err != nil {
		return
	}
	t.Log(re)
}

func TestMongoFindAllData(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%d",
		cfg.MongoCfg.Host, cfg.MongoCfg.Port,
	)
	mongoService, err := dbservice.NewMongoDBClient(mongoLink, cfg.MongoCfg.DB)
	if err != nil {
		panic("Failed to initialize Mongo: " + err.Error())
	}
	filter := bson.M{
		"age": bson.M{"$gt": 29}, // $gt 代表大于
	}
	re, err := mongoService.FindAll("go_test", filter)
	if err != nil {
		return
	}
	t.Log(re)
}

func TestMongoDropCollection(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	mongoLink := fmt.Sprintf(
		"mongodb://%s:%d",
		cfg.MongoCfg.Host, cfg.MongoCfg.Port,
	)
	mongoService, err := dbservice.NewMongoDBClient(mongoLink, cfg.MongoCfg.DB)
	if err != nil {
		panic("Failed to initialize Mongo: " + err.Error())
	}
	err = mongoService.DropCollection("go_test")
	if err != nil {
		return
	}
}
