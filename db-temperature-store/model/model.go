package model

import (
	"fmt"
	"gorm.io/gorm"
)

type TemperatureUpdate struct {
	gorm.Model         // Not necessary but brings a pack of extra fields managed by convention (ID, timestamps, etc)
	SensorId   string  `json:"sensorId"`
	Value      float64 `json:"value"`
}

func (t TemperatureUpdate) String() string {
	return fmt.Sprintf("{sensorId=%s, value=%f}", t.SensorId, t.Value)
}
