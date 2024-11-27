package data_flow_model

import "time"

type RecAircraftStatusRequest struct {
	AircraftID int `json:"AircraftID"`
}

type TaskAircraftRedis struct {
	TaskID     int        `json:"TaskID"`
	AircraftID int        `json:"AircraftID"`
	LaneID     int        `json:"LaneID"`
	CreateTime time.Time  `json:"CreateTime"`
	EndTime    *time.Time `json:"EndTime"`
	TrackTable string     `json:"TrackTable"`
	EventTable string     `json:"EventTable"`
	TimeStr    string     `json:"TimeStr"`
}
