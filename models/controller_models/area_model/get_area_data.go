package area_model

// GetAreaData 结构体表示获取区域数据的请求
type GetAreaData struct {
	AreaID int `json:"AreaID"` // 区域ID
}

type GetRasterData struct {
	AreaID int    `json:"AreaID"` // 区域ID
	Status string `json:"Status"` // 区域ID
}

type GetRasterSizeMongo struct {
	RasterSize []interface{} `json:"RasterSize"` // 栅格大小
}
