package model

import (
	"math"
	"slices"
	"sync"
)

type SyncSensorsList struct {
	Mu   sync.Mutex
	Wg   sync.WaitGroup
	List []Sensor
}

func (s *SyncSensorsList) GetSensor() *Sensor {
	if len(s.List) == 0 {
		return nil
	}
	s.SortByAqi()
	trustedIndex := math.Ceil(float64(len(s.List) / 2))
	return &s.List[int(trustedIndex)]
}

func (s *SyncSensorsList) AddSensor(sensor Sensor) {
	s.Mu.Lock()
	s.List = append(s.List, sensor)
	s.Mu.Unlock()
}

func (s *SyncSensorsList) SortByAqi() {
	slices.SortFunc(s.List, func(a, b Sensor) int {
		if a.Aqi < b.Aqi {
			return -1
		} else if a.Aqi > b.Aqi {
			return 1
		}
		return 0
	})
}
