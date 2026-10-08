package dto

type Response struct {
	Answers Answers `json:"answers"`
}

type Answers struct {
	Position PositionResponse `json:"position"`
}

type PositionResponse struct {
	Choice        string             `json:"choice"`
	Legend        map[string]string  `json:"legend"`
	Probabilities map[string]float32 `json:"probabilities"`
}
