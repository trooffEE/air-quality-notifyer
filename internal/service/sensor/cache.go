package sensor

import (
	rSensor "air-quality-notifyer/internal/db/repository/sensor"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

var TTL = time.Hour * 4

func (s *Service) saveSensorInCache(ctx context.Context, sensor rSensor.Sensor) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
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

	if _, err = pipeline.Exec(ctx); err != nil {
		zap.L().Error("failed to execute sensor cache pipeline", zap.Error(err), zap.Any("payload", payload))
	}
}

func (s *Service) getSensorFromCache(ctx context.Context, sensorId int64) (*rSensor.Sensor, error) {
	key := getSensorCacheKey(sensorId)

	result, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var sensor rSensor.Sensor
	err = json.Unmarshal([]byte(result), &sensor)
	if err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (s *Service) getDistrictSensorsFromCache(ctx context.Context, districtID int64) (*[]rSensor.Sensor, error) {
	key := getDistrictSensorsCacheKey(districtID)

	result, err := s.cache.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var sensors []rSensor.Sensor
	for _, sensorJSON := range result {
		var sensor rSensor.Sensor
		if err := json.Unmarshal([]byte(sensorJSON), &sensor); err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}

	return &sensors, nil
}

func (s *Service) GetAliveSensorsFromCache(ctx context.Context) ([]AliveSensor, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var cursor uint64
	var payloads []string
	for {
		keys, nextCursor, err := s.cache.Scan(ctx, cursor, "sensor:*", 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			if !isAliveSensorCacheKey(key) {
				continue
			}

			payload, err := s.cache.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			payloads = append(payloads, payload)
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return aliveSensorsFromCachePayloads(payloads)
}

func aliveSensorsFromCachePayloads(payloads []string) ([]AliveSensor, error) {
	sensors := make([]AliveSensor, 0, len(payloads))
	for _, payload := range payloads {
		var sensor rSensor.Sensor
		if err := json.Unmarshal([]byte(payload), &sensor); err != nil {
			return nil, err
		}

		sensors = append(sensors, AliveSensor{
			APIID:    sensor.ApiId,
			Address:  sensor.Address,
			Lat:      sensor.Lat,
			Lon:      sensor.Lon,
			District: sensor.District.Name,
		})
	}

	return sensors, nil
}

func isAliveSensorCacheKey(key string) bool {
	return strings.HasPrefix(key, "sensor:") && !strings.HasPrefix(key, "sensor:district:")
}

func getSensorCacheKey(sensorId int64) string {
	return fmt.Sprintf("sensor:%d", sensorId)
}

func getDistrictSensorsCacheKey(districtId int64) string {
	return fmt.Sprintf("sensor:district:%d", districtId)
}
