package lane_model

// LaneCreateModel 结构体表示创建航线的模型
type LaneCreateModel struct {
	Name       string      `json:"Name"`       // 航线名称
	IsHide     bool        `json:"IsHide"`     // 是否隐藏
	PointData  [][]float64 `json:"PointData"`  // 点数据
	RasterData []uint      `json:"RasterData"` // 栅格数据
}
