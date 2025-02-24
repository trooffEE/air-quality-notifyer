package sensor

import (
	"go.uber.org/zap"
	"slices"
)

func (s *Service) invalidateSensors(aliveSensors []AqiSensorScriptScrapped) {
	currentlySavedSensorsIds, err := s.repo.GetAllApiIds()
	if err != nil {
		zap.L().Error("failed to get airquality sensors from db", zap.Error(err))
		return
	}

	var aliveIds []int64
	for _, sensor := range aliveSensors {
		aliveIds = append(aliveIds, sensor.Id)
	}

	for _, id := range currentlySavedSensorsIds {
		if !slices.Contains(aliveIds, id) {
			err := s.repo.EvictSensor(id)
			if err != nil {
				zap.L().Error("failed to evict sensor from db", zap.Int64("id", id), zap.Error(err))
			}
		}
	}
}
