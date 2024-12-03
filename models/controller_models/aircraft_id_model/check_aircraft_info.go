package aircraft_id_model

import "time"

// GetAircraftInfoID 用于获取飞机信息的结构体
type GetAircraftInfoID struct {
	AircraftID int `json:"AircraftID"` // 飞机ID
}

// SetAircraftInfo 用于设置飞机信息的结构体
type SetAircraftInfo struct {
	Company string `json:"Company"` // 公司
	Name    string `json:"Name"`    // 名称
	Type    string `json:"Type"`    // 类型
}

// MysqlAircraftInfo 用于表示飞机信息的结构体
type MysqlAircraftInfo struct {
	Company    string    `json:"Company"`    // 公司
	Name       string    `json:"Name"`       // 名称
	Type       string    `json:"Type"`       // 类型
	TimeStr    string    `json:"TimeStr"`    // 时间字符串
	CreateTime time.Time `json:"CreateTime"` // 创建时间
	AircraftID int       `json:"AircraftID"` // 飞机ID
}
