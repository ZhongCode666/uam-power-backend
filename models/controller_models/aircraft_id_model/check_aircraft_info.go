package aircraft_id_model

import "time"

type GetAircraftInfoID struct {
	AircraftID int `json:"AircraftID"`
}

type SetAircraftInfo struct {
	Company string `json:"Company"`
	Name    string `json:"Name"`
	Type    string `json:"Type"`
}

type MysqlAircraftInfo struct {
	Company    string    `json:"Company"`
	Name       string    `json:"Name"`
	Type       string    `json:"Type"`
	TimeStr    string    `json:"TimeStr"`
	CreateTime time.Time `json:"CreateTime"`
	AircraftID int       `json:"AircraftID"`
}
