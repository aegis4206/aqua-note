package models

type SensorData struct {
	Device      string  `json:"device"`
	Temperature float64 `json:"temperature"`
	TdsPpm      float64 `json:"tds_ppm,omitempty"`
}
