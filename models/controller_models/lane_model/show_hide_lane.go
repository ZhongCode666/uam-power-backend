package lane_model

// ShowHideLaneModel 结构体表示显示或隐藏航线的模型
type ShowHideLaneModel struct {
	// LaneIDs 是航线ID的列表
	LaneIDs []int `json:"LaneIDs"`
}
