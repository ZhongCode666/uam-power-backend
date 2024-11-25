package area_model

import "time"

type AreaMysqlModel struct {
	AreaID     int       `json:"AreaID"`
	CreateTime time.Time `json:"CreateTime"`
	Name       string    `json:"Name"`
	TimeStr    string    `json:"TimeStr"`
}
