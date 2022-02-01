package model

import "fmt"

type TemperatureUpdate struct {
	SensorId string  `json:"sensorId"`
	Value    float64 `json:"value"`
}

func (t TemperatureUpdate) String() string {
	return fmt.Sprintf("{sensorId=%s, value=%f}", t.SensorId, t.Value)
}
