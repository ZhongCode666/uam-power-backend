package area_model

// RasterData 表示栅格数据的结构体
type RasterData struct {
	AreaID    int   `json:"AreaID"`    // 区域ID
	RasterIDs []int `json:"RasterIDs"` // 栅格ID列表
}
