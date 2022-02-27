package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
)

type TemperatureUpdate struct {
	gorm.Model         // Not necessary but brings a pack of extra fields managed by convention (ID, timestamps, etc)
	SensorId   string  `json:"sensorId"`
	Value      float64 `json:"value"`
}

func (t TemperatureUpdate) String() string {
	return fmt.Sprintf("{sensorId=%s, value=%f}", t.SensorId, t.Value)
}

func TemperatureUpdateFromJSONBytes(jsonBytes []byte) (*TemperatureUpdate, error) {
	jsonObject := &TemperatureUpdate{}
	err := json.Unmarshal(jsonBytes, jsonObject)
	return jsonObject, err
}

func TemperatureUpdateFromJSONReader(reader io.Reader) (*TemperatureUpdate, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return TemperatureUpdateFromJSONBytes(data)
}

func (t *TemperatureUpdate) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TemperatureUpdate) ToJSONReader() (io.Reader, error) {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBytes), nil
}
