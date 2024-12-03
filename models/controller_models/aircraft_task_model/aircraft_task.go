package aircraft_task_model

import "time"

// CreateTaskAircraftInfo 用于创建任务的飞机信息结构体
type CreateTaskAircraftInfo struct {
	AircraftID int `json:"AircraftID"` // 飞机ID
	LaneID     int `json:"LaneID"`     // 航线ID
}

// MysqlAircraftTask 用于表示任务的飞机信息结构体
type MysqlAircraftTask struct {
	TaskID     int        `json:"TaskID"`     // 任务ID
	AircraftID int        `json:"AircraftID"` // 飞机ID
	LaneID     int        `json:"LaneID"`     // 航线ID
	CreateTime time.Time  `json:"CreateTime"` // 创建时间
	EndTime    *time.Time `json:"EndTime"`    // 结束时间
	TrackTable string     `json:"TrackTable"` // 轨迹表
	EventTable string     `json:"EventTable"` // 事件表
	TimeStr    string     `json:"TimeStr"`    // 时间字符串
}

// ByAircraftID 用于根据飞机ID查询的结构体
type ByAircraftID struct {
	AircraftID int `json:"AircraftID"` // 飞机ID
}

// ByAircraftIDAndTaskID 用于根据飞机ID和任务ID查询的结构体
type ByAircraftIDAndTaskID struct {
	AircraftID int `json:"AircraftID"` // 飞机ID
	TaskID     int `json:"TaskID"`     // 任务ID
}
