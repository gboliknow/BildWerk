package processing

type TransformRequest struct {
	Transformations struct {
		Resize struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"resize"`
		Crop struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			X      int `json:"x"`
			Y      int `json:"y"`
		} `json:"crop"`
		Rotate    int    `json:"rotate"`
		Format    string `json:"format"`
		Compress  bool   `json:"compress"`
		Flip      bool   `json:"flip"`
		Mirror    bool   `json:"mirror"`
		Watermark struct {
			Text     string `json:"text"`
			Position string `json:"position"`
		} `json:"watermark"`
		Filters struct {
			Grayscale bool `json:"grayscale"`
			Sepia     bool `json:"sepia"`
		} `json:"filters"`
	} `json:"transformations"`
}

type UploadImageDTO struct {
	Filename string `form:"filename" binding:"required"`
}
