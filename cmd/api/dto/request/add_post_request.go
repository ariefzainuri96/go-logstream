package request

import (
	"encoding/json"
)

type AddPostRequest struct {
	ProjectId uint   `json:"project_id" validate:"required"`
	Title     string `json:"title" validate:"required,max=255"`
	Content   string `json:"content" validate:"required"`
	Category  string `json:"category" validate:"max=50"`
	Status    string `json:"status" validate:"max=20"`
}

func (r AddPostRequest) Marshal() ([]byte, error) {
	marshal, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return marshal, nil
}

func (r *AddPostRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}
