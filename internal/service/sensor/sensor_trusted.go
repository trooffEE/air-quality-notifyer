package sensor

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func (s *Service) GetTrustedSensorsEveryHour() {
	cronCreator := cron.New()
	cronString := "0 * * * *"

	_, err := cronCreator.AddFunc(cronString, func() {
		if time.Now().UTC().Hour()%AliveSensorTimeDiff == 0 {
			<-s.syncCron
		}
		s.getTrustedAirqualitySensors()
	})
	if err != nil {
		panic(err)
	}

	cronCreator.Start()
}

func (s *Service) getTrustedAirqualitySensors() {
	allDistricts := s.sDistricts.GetAllDistricts() // think about it

	respChan := make(chan AqiSensor, len(allDistricts))
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

	var sensors []AqiSensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	s.cSensors <- sensors
}
