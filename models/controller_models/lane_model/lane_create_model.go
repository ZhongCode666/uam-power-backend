package lane_model

type LaneCreateModel struct {
	Name       string      `json:"Name"`
	IsHide     bool        `json:"IsHide"`
	PointData  [][]float64 `json:"PointData"`
	RasterData []uint      `json:"RasterData"`
}
