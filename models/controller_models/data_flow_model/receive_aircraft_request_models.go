package data_flow_model

import "time"

// RecAircraftStatusRequest 结构体表示接收飞机状态的请求
type RecAircraftStatusRequest struct {
	AircraftID int `json:"AircraftID"` // 飞机ID
}

// TaskAircraftRedis 结构体表示任务航线的 Redis 数据
type TaskAircraftRedis struct {
	TaskID     int        `json:"TaskID"`     // 任务ID
	AircraftID int        `json:"AircraftID"` // 飞机ID
	LaneID     int        `json:"LaneID"`     // 航线ID
	CreateTime time.Time  `json:"CreateTime"` // 创建时间
	EndTime    *time.Time `json:"EndTime"`    // 结束时间
	TrackTable string     `json:"TrackTable"` // 轨迹表
	EventTable string     `json:"EventTable"` // 事件表
	TimeStr    string     `json:"TimeStr"`    // 时间字符串
}
