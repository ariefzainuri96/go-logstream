package request

import (
	"encoding/json"
)

type CheckSlugRequest struct {
	Slug string `json:"slug" validate:"required,max=255"`
}

func (r CheckSlugRequest) Marshal() ([]byte, error) {
	marshal, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return marshal, nil
}

func (r *CheckSlugRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}
