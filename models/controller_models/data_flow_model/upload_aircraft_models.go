package data_flow_model

// AircraftStatus 结构体表示飞机状态
type AircraftStatus struct {
	TimeString string  `json:"TimeString"` // 时间字符串
	Yaw        float64 `json:"Yaw"`        // 偏航角
	Latitude   float64 `json:"Latitude"`   // 纬度
	Longitude  float64 `json:"Longitude"`  // 经度
	Altitude   float64 `json:"Altitude"`   // 高度
	AircraftID int     `json:"AircraftID"` // 飞机ID
}

// AircraftEvent 结构体表示飞机事件
type AircraftEvent struct {
	TimeString string `json:"TimeString"` // 时间字符串
	Event      string `json:"Event"`      // 事件
	AircraftID int    `json:"AircraftID"` // 飞机ID
}
