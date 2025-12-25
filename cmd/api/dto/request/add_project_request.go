package request

import (
	"encoding/json"
)

type AddProjectRequest struct {
	Name       string `json:"name" validate:"required,max=255"`
	Slug       string `json:"slug" validate:"required,max=255"`
	WebhookUrl string `json:"webhook_url"`
}

func (r AddProjectRequest) Marshal() ([]byte, error) {
	marshal, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return marshal, nil
}

func (r *AddProjectRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &r)
}
