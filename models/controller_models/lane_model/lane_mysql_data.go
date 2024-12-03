package lane_model

import "time"

// LaneMysqlData 结构体表示航线的 MySQL 数据
type LaneMysqlData struct {
	LaneID     int       `json:"LaneID"`     // 航线ID
	CreateTime time.Time `json:"CreateTime"` // 创建时间
	Name       string    `json:"Name"`       // 航线名称
	IsHide     int       `json:"IsHide"`     // 是否隐藏
	TimeStr    string    `json:"TimeStr"`    // 时间字符串
}
