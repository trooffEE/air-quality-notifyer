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

type SyncSensors struct {
	mu   sync.Mutex
	wg   sync.WaitGroup
	list []Sensor
}

// "Trusted" - meaning median from what we have in district, slightly more realistic than AVG and Worst AQI value determination
func (s *SyncSensors) getTrustedSensor() *Sensor {
	if len(s.list) == 0 {
		return nil
	}
	s.sortByAqi()
	trustedIndex := math.Ceil(float64(len(s.list) / 2))
	return &s.list[int(trustedIndex)]
}

func (s *SyncSensors) sortByAqi() {
	slices.SortFunc(s.list, func(a, b Sensor) int {
		if a.Aqi < b.Aqi {
			return -1
		} else if a.Aqi > b.Aqi {
			return 1
		}
		return 0
	})
}

func (s *SyncSensors) addSensor(sensor Sensor) {
	s.mu.Lock()
	s.list = append(s.list, sensor)
	s.mu.Unlock()
}

func (s *Service) GetTrustedSensorsEveryHour() {
	cronCreator := cron.New()
	cronString := "0 * * * *"

	_, err := cronCreator.AddFunc(cronString, func() {
		if time.Now().UTC().Hour()%AliveSensorTimeDiff == 0 {
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
	wg.Add(len(allDistricts))
	for _, district := range allDistricts {
		sensorsInDistrict, err := s.getDistrictSensorsFromCache(district.Id)
		if err != nil || sensorsInDistrict == nil {
			zap.L().Error("failed to get sensors by districtId", zap.Error(err), zap.Int64("districtId", district.Id))
			continue
		}
		go func() {
			defer wg.Done()
			findTrustedSensor(respChan, *sensorsInDistrict)
		}()
	}
	wg.Wait()
	close(respChan)

	var sensors []Sensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	s.cSensors <- sensors
}

func findTrustedSensor(resChan chan Sensor, sensors []models.AirqualitySensor) {
	var syncSensorList SyncSensors
	syncSensorList.wg.Add(len(sensors))
	for _, sensor := range sensors {
		go getLastUpdatedSensor(&syncSensorList, sensor.ApiId, sensor.District.Name)
	}
	syncSensorList.wg.Wait()

	trustedAqlSensor := syncSensorList.getTrustedSensor()
	if trustedAqlSensor == nil {
		return
	}

	resChan <- *trustedAqlSensor
}
