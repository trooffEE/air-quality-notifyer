package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var TTL time.Duration = time.Hour * 4

func (s *Service) saveSensorInCache(sensor models.AirqualitySensor) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pipeline := s.cache.TxPipeline()

	key := getSensorCacheKey(sensor.ApiId)
	payload, err := json.Marshal(sensor)
	if err != nil {
		zap.L().Error("failed to marshal sensor", zap.Error(err), zap.Any("payload", payload))
		return
	}

	err = pipeline.Set(
		ctx,
		key,
		payload,
		TTL,
	).Err()

	if err != nil {
		zap.L().Error("failed to save sensor", zap.Error(err), zap.Any("payload", payload))
	}

	districtKey := getDistrictSensorsCacheKey(sensor.DistrictId)

	err = pipeline.HSet(ctx, districtKey, key, payload).Err()

	if err != nil {
		zap.L().Error(
			fmt.Sprintf("failed to update hash set of %s", districtKey),
			zap.Error(err),
			zap.Any("payload", payload),
		)
	}

	err = pipeline.Expire(ctx, districtKey, TTL).Err()

	if err != nil {
		zap.L().Error(
			fmt.Sprintf("failed to set expiration of %s", districtKey),
			zap.Error(err),
			zap.Any("payload", payload),
		)
	}

	pipeline.Exec(ctx)
}

func (s *Service) getSensorFromCache(sensorId int64) (*models.AirqualitySensor, error) {
	key := getSensorCacheKey(sensorId)

	result, err := s.cache.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var sensor models.AirqualitySensor
	err = json.Unmarshal([]byte(result), &sensor)
	if err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (s *Service) getDistrictSensorsFromCache(districtID int64) (*[]models.AirqualitySensor, error) {
	key := getDistrictSensorsCacheKey(districtID)

	result, err := s.cache.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var sensors []models.AirqualitySensor
	for _, sensorJSON := range result {
		var sensor models.AirqualitySensor
		if err := json.Unmarshal([]byte(sensorJSON), &sensor); err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}

	return &sensors, nil
}

func getSensorCacheKey(sensorId int64) string {
	return fmt.Sprintf("sensor:%d", sensorId)
}

func getDistrictSensorsCacheKey(districtId int64) string {
	return fmt.Sprintf("sensor:district:%d", districtId)
}
