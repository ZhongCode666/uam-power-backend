package aircraft_task_model

import "time"

type CreateTaskAircraftInfo struct {
	AircraftID int `json:"AircraftID"`
	LaneID     int `json:"LaneID"`
}

type MysqlAircraftTask struct {
	TaskID     int        `json:"TaskID"`
	AircraftID int        `json:"AircraftID"`
	LaneID     int        `json:"LaneID"`
	CreateTime time.Time  `json:"CreateTime"`
	EndTime    *time.Time `json:"EndTime"`
	TrackTable string     `json:"TrackTable"`
	EventTable string     `json:"EventTable"`
	TimeStr    string     `json:"TimeStr"`
}

type ByAircraftID struct {
	AircraftID int `json:"AircraftID"`
}

type ByAircraftIDAndTaskID struct {
	AircraftID int `json:"AircraftID"`
	TaskID     int `json:"TaskID"`
}
