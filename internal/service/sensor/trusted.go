package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

/**
"Trusted" is just median AQI in district
*/

func (s *Service) StartGettingTrustedSensorsEveryHour() {
	cronCreator := cron.New()
	cronString := "0 * * * *"

	_, err := cronCreator.AddFunc(cronString, func() {
		if time.Now().UTC().Hour()%InvalidationPeriod == 0 {
			<-s.syncCron
		}
		s.getTrustedSensors()
	})
	if err != nil {
		panic(err)
	}

	cronCreator.Start()
}

func (s *Service) getTrustedSensors() {
	allDistricts := s.sDistricts.GetAllDistricts() // think about it

	respChan := make(chan Sensor, len(allDistricts))
	wg := sync.WaitGroup{}
	for _, district := range allDistricts {
		sensorsInDistrict, err := s.getDistrictSensorsFromCache(district.Id)
		if err != nil || sensorsInDistrict == nil {
			zap.L().Error("failed to get sensors by districtId", zap.Error(err), zap.Int64("districtId", district.Id))
			continue
		}
		wg.Go(func() { findTrustedSensor(respChan, *sensorsInDistrict) })
	}
	wg.Wait()
	close(respChan)

	var sensors []Sensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	s.cSensors <- sensors
}

func findTrustedSensor(resChan chan Sensor, sensors []models.Sensor) {
	var syncSensorList SyncTrustedSensors
	syncSensorList.wg.Add(len(sensors))
	for _, sensor := range sensors {
		go getLastUpdatedSensor(&syncSensorList, sensor.ApiId, sensor.District.Name)
	}
	syncSensorList.wg.Wait()

	trustedAqlSensor := syncSensorList.getSensor()
	if trustedAqlSensor == nil {
		return
	}

	resChan <- *trustedAqlSensor
}

type SyncTrustedSensors struct {
	mu   sync.Mutex
	wg   sync.WaitGroup
	list []Sensor
}

func (s *SyncTrustedSensors) getSensor() *Sensor {
	if len(s.list) == 0 {
		return nil
	}
	s.sortByAqi()
	trustedIndex := math.Ceil(float64(len(s.list) / 2))
	return &s.list[int(trustedIndex)]
}

func (s *SyncTrustedSensors) addSensor(sensor Sensor) {
	s.mu.Lock()
	s.list = append(s.list, sensor)
	s.mu.Unlock()
}

func (s *SyncTrustedSensors) sortByAqi() {
	slices.SortFunc(s.list, func(a, b Sensor) int {
		if a.Aqi < b.Aqi {
			return -1
		} else if a.Aqi > b.Aqi {
			return 1
		}
		return 0
	})
}
