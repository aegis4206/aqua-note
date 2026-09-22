package model

type SensorData struct {
	DeviceCode  *string  `json:"device_code"`
	Temperature *float64 `json:"temperature"`
}
