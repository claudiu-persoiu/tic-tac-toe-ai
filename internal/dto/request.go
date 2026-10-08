package dto

type Request struct {
	Model     string            `json:"model"`
	State     map[string]string `json:"state"`
	Questions Questions         `json:"questions"`
}

type Questions struct {
	Position PositionRequest `json:"position"`
}

type PositionRequest struct {
	Type         string            `json:"type"`
	Instructions string            `json:"instructions"`
	Criteria     map[string]string `json:"criteria"`
}
