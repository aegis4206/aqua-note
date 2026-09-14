package models

type SensorData struct {
	Device_code string  `json:"device_code"`
	Temperature float64 `json:"temperature"`
	TdsPpm      float64 `json:"tds_ppm,omitempty"`
}
