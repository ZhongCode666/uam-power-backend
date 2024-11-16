package data_flow_model

type AircraftStatus struct {
	TimeString string  `json:"TimeString"`
	Yaw        float64 `json:"Yaw"`
	Latitude   float64 `json:"Latitude"`
	Longitude  float64 `json:"Longitude"`
	Altitude   float64 `json:"Altitude"`
	AircraftID int     `json:"AircraftID"`
}

type AircraftEvent struct {
	TimeString string `json:"TimeString"`
	Event      string `json:"Event"`
	AircraftID int    `json:"AircraftID"`
}
