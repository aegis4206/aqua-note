package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"aquanote-backend/internal/model"
)

type TemperatureLogRepository struct {
	db *gorm.DB
}

func NewTemperatureLogRepository(db *gorm.DB) *TemperatureLogRepository {
	return &TemperatureLogRepository{db: db}
}

func (r *TemperatureLogRepository) Create(ctx context.Context, data model.SensorData) (*model.TemperatureLog, error) {
	logEntry := model.TemperatureLog{
		DeviceCode:  data.DeviceCode,
		Temperature: data.Temperature,
	}

	if err := r.db.WithContext(ctx).Create(&logEntry).Error; err != nil {
		return nil, fmt.Errorf("insert temperature log failed: %w", err)
	}

	return &logEntry, nil
}
