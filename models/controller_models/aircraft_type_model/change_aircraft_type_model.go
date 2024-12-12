package aircraft_type_model

type ChangeAircraftTypeModel struct {
	Type   string  `json:"Type"`   // 飞机类型
	GLB    string  `json:"GLB"`    // GLB 文件
	ScaleX float32 `json:"ScaleX"` // 缩放比例
	ScaleY float32 `json:"ScaleY"` // 缩放比例
	ScaleZ float32 `json:"ScaleZ"` // 缩放比例
}
