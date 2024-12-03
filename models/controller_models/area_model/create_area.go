package area_model

// CreateArea 用于创建区域的结构体
type CreateArea struct {
	Name        string                   `json:"Name"`        // 区域名称
	RangeData   []float64                `json:"RangeData"`   // 范围数据
	RasterSize  []float64                `json:"RasterSize"`  // 栅格大小
	RasterIndex [][][]int                `json:"RasterIndex"` // 栅格索引
	RasterData  map[int]SingleRasterData `json:"RasterData"`  // 栅格数据
}

// SingleRasterData 表示单个栅格数据的结构体
type SingleRasterData struct {
	Status    string  `json:"Status"`    // 状态
	Longitude float64 `json:"Longitude"` // 经度
	Latitude  float64 `json:"Latitude"`  // 纬度
	Altitude  float64 `json:"Altitude"`  // 高度
}

// GenerateAreaID 用于生成区域 ID 的结构体
type GenerateAreaID struct {
	Name string `json:"Name"` // 区域名称
}

// UploadArea 用于上传区域的结构体
type UploadArea struct {
	AreaID      int                      `json:"AreaID"`      // 区域ID
	RangeData   []float64                `json:"RangeData"`   // 范围数据
	RasterSize  []float64                `json:"RasterSize"`  // 栅格大小
	RasterIndex [][][]int                `json:"RasterIndex"` // 栅格索引
	RasterData  map[int]SingleRasterData `json:"RasterData"`  // 栅格数据
}
