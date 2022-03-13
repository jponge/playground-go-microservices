package data

import (
	"encoding/json"
	"io"
)

type Payload struct {
	SensorID string  `json:"sensorId"`
	Value    float64 `json:"value"`
}

func PayloadFromBytes(jsonBytes []byte) (*Payload, error) {
	payload := &Payload{}
	err := json.Unmarshal(jsonBytes, payload)
	return payload, err
}

func PayloadFromReader(reader io.Reader) (*Payload, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return PayloadFromBytes(data)
}

func (payload Payload) ToJSON() ([]byte, error) {
	return json.Marshal(payload)
}
