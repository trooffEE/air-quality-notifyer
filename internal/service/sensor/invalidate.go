package sensor

import (
	"fmt"
	"slices"
)

func (s *Service) invalidateSensors(aliveSensors []AqiSensorScriptScrapped) {
	currentlySavedSensorsIds, err := s.repo.GetAllApiIds()
	if err != nil {
		fmt.Printf("Failed to get all air quality sensor: %+v\n", err)
		return
	}

	var aliveIds []int64
	for _, sensor := range aliveSensors {
		aliveIds = append(aliveIds, sensor.Id)
	}

	for _, id := range *currentlySavedSensorsIds {
		if !slices.Contains(aliveIds, id) {
			err := s.repo.EvictSensor(id)
			if err != nil {
				fmt.Printf("Failed to evict sensor: %+v\n", err)
			}
		}
	}
}
