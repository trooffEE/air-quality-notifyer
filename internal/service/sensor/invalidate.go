package sensor

import (
	"air-quality-notifyer/internal/lib"
	"slices"
)

func (s *Service) invalidateSensors(aliveSensors []AqiSensorScriptScrapped) {
	currentlySavedSensorsIds, err := s.repo.GetAllApiIds()
	if err != nil {
		lib.LogError("invalidateSensors", "failed to get airquality sensors from db", err)
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
				lib.LogError("invalidateSensors", "failed to evict sensor from db", err)
			}
		}
	}
}
