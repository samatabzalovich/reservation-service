package sms

import (
	"encoding/json"
)

func (service *MessageService) toJson(data any) ([]byte, error) {
	out, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return out, nil
}
