package area_model

import "time"

// AreaMysqlModel 用于表示区域信息的结构体
type AreaMysqlModel struct {
	AreaID     int       `json:"AreaID"`     // 区域ID
	CreateTime time.Time `json:"CreateTime"` // 创建时间
	Name       string    `json:"Name"`       // 名称
	TimeStr    string    `json:"TimeStr"`    // 时间字符串
}
