package model

import (
	"sync"
	"time"
)

type TemperatureLog struct {
	ID          uint64    `gorm:"column:id;primaryKey" json:"id"`
	DeviceCode  *string   `gorm:"column:device_code" json:"device_code,omitempty"`
	Temperature *float64  `gorm:"column:temperature" json:"temperature,omitempty"`
	RecordedAt  time.Time `gorm:"column:recorded_at;autoCreateTime" json:"recorded_at"`
}

func (TemperatureLog) TableName() string {
	return "temperature_logs"
}

// TCSDeviceMapping TCS 設備資訊
type TCSDeviceMapping struct {
	ID         uint64   `gorm:"column:id;primaryKey" json:"id"`
	DeviceCode string   `gorm:"column:device_code;not null" json:"device_code"`
	DeviceName string   `gorm:"column:device_name;not null" json:"device_name"`
	ZoneID     uint64   `gorm:"column:zone_id;not null" json:"zone_id"`
	ZoneName   string   `gorm:"column:zone_name;not null" json:"zone_name"`
	UpperLimit *float64 `gorm:"column:upper_limit" json:"upper_limit,omitempty"`
	LowerLimit *float64 `gorm:"column:lower_limit" json:"lower_limit,omitempty"`
	Memo       *string  `gorm:"column:memo" json:"memo,omitempty"`
	// ModifierID *uint64        `gorm:"column:modifier_id" json:"modifier_id,omitempty"`
	// CreatorID  *uint64        `gorm:"column:creator_id" json:"creator_id,omitempty"`
	// CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at,omitempty"`
	// UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at,omitempty"`
	// DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (TCSDeviceMapping) TableName() string {
	return "device_mappings"
}

type ExceedLimitTracker struct {
	Lock    sync.Mutex
	Devices map[string]int
}

func (e *ExceedLimitTracker) Update(deviceCode string, exceeded bool) int {
	e.Lock.Lock()
	defer e.Lock.Unlock()

	if exceeded {
		e.Devices[deviceCode]++
	} else {
		e.Devices[deviceCode] = 0
	}

	return e.Devices[deviceCode]
}

type NotifyRequest struct {
	Type               string `json:"type"` // up or low
	DeviceCode         string `json:"device_code"`
	DeviceName         string `json:"device_name"`
	ZoneName           string `json:"zone_name"`
	CurrentTemperature string `json:"current_temperature"`
	Limit              string `json:"limit"`
}
