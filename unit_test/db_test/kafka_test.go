package db

import (
	"context"
	"testing"
	"time"
	"uam-power-backend/service/db_service"
	"uam-power-backend/unit_test/config_test"
)

func TestKafkaConsumerRec(t *testing.T) {
	cfg := DBconfig.NewConfig()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// 初始化服务
	KafkaConsumerService := dbservice.NewKafkaConsumer(cfg.KafkaAddr, cfg.KafkaTopic, "1")
	re, _ := KafkaConsumerService.ReceiveMessage(ctx)
	t.Log(re)
	err := KafkaConsumerService.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestKafkaProducerPro(t *testing.T) {
	cfg := DBconfig.NewConfig()
	// 初始化服务
	KafkaProducerService := dbservice.NewKafkaProducer(cfg.KafkaAddr, cfg.KafkaTopic)
	t.Log(KafkaProducerService.SendMessage("go_test"))
	err := KafkaProducerService.Close()
	if err != nil {
		t.Error(err)
	}
}
