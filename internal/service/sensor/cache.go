package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var TTL time.Duration = time.Hour * 4

func (s *Service) saveSensorInCache(sensor models.AirqualitySensor) {
	payload, err := json.Marshal(sensor)
	if err != nil {
		zap.L().Error("failed to marshal sensor", zap.Error(err), zap.Any("payload", payload))
		return
	}

	status := s.cache.Set(
		context.Background(),
		getSensorCacheKey(sensor.ApiId),
		payload,
		TTL,
	)

	if err := status.Err(); status.Err() != nil {
		zap.L().Error("failed to save sensor", zap.Error(err), zap.Any("payload", payload))
	}
}

func (s *Service) getSensorFromCache(sensorId int64) (*models.AirqualitySensor, error) {
	result, err := s.cache.Get(context.Background(), getSensorCacheKey(sensorId)).Result()
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

func getSensorCacheKey[K int64 | string](sensorId K) string {
	var cacheKeyPrefix = "sensor:"
	switch id := any(sensorId).(type) {
	case string:
		return cacheKeyPrefix + id
	case int64:
		return cacheKeyPrefix + strconv.Itoa(int(id))
	default:
		return cacheKeyPrefix
	}
}
