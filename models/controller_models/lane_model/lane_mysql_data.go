package lane_model

import "time"

type LaneMysqlData struct {
	LaneID     int       `json:"LaneID"`
	CreateTime time.Time `json:"CreateTime"`
	Name       string    `json:"Name"`
	IsHide     int       `json:"IsHide"`
	TimeStr    string    `json:"TimeStr"`
}
