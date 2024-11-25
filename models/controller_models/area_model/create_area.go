package area_model

type CreateArea struct {
	Name        string                   `json:"Name"`
	RangeData   []float64                `json:"RangeData"`
	RasterSize  []float64                `json:"RasterSize"`
	RasterIndex [][][]int                `json:"RasterIndex"`
	RasterData  map[int]SingleRasterData `json:"RasterData"`
}

type SingleRasterData struct {
	Status    string  `json:"Status"`
	Longitude float64 `json:"Longitude"`
	Latitude  float64 `json:"Latitude"`
	Altitude  float64 `json:"Altitude"`
}
