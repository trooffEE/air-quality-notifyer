package sensor

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"air-quality-notifyer/internal/service/sensor/model"
	"air-quality-notifyer/internal/service/sensor/request"
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

	respChan := make(chan model.Sensor, len(allDistricts))
	wg := sync.WaitGroup{}
	for _, district := range allDistricts {
		sensorsInDistrict, err := s.getDistrictSensorsFromCache(district.Id)
		if err != nil || sensorsInDistrict == nil {
			zap.L().Error("failed to get sensors by districtId", zap.Error(err), zap.Int64("districtId", district.Id))
			continue
		}
		wg.Go(func() { getTrustedSensor(respChan, *sensorsInDistrict) })
	}
	wg.Wait()
	close(respChan)

	var sensors []model.Sensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	s.cSensors <- sensors
}

func getTrustedSensor(resChan chan model.Sensor, sensors []sensor.Sensor) {
	var syncSensorList model.SyncSensorsList
	syncSensorList.Wg.Add(len(sensors))
	for _, sensor := range sensors {
		go request.GetArchiveSensor(&syncSensorList, sensor.ApiId, sensor.District.Name)
	}
	syncSensorList.Wg.Wait()

	trustedAqlSensor := syncSensorList.GetSensor()
	if trustedAqlSensor == nil {
		return
	}

	resChan <- *trustedAqlSensor
}
