package sensor

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"air-quality-notifyer/internal/service/sensor/model"
	"air-quality-notifyer/internal/service/sensor/request"
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

/**
"Trusted" is just median AQI in district
*/

func (s *Service) StartGettingTrustedSensorsEveryHour(ctx context.Context) func(context.Context) {
	cronCreator := cron.New()
	cronString := "0 * * * *"

	_, err := cronCreator.AddFunc(cronString, func() {
		if time.Now().UTC().Hour()%InvalidationPeriod == 0 {
			select {
			case <-s.syncCron:
			case <-ctx.Done():
				return
			}
		}
		s.getTrustedSensors(ctx)
	})
	if err != nil {
		panic(err)
	}

	cronCreator.Start()

	return func(shutdownCtx context.Context) {
		stopCron(shutdownCtx, cronCreator)
	}
}

func (s *Service) getTrustedSensors(ctx context.Context) {
	allDistricts := s.sDistricts.GetAllDistricts(ctx) // think about it

	respChan := make(chan model.Sensor, len(allDistricts))
	wg := sync.WaitGroup{}
	for _, district := range allDistricts {
		if ctx.Err() != nil {
			return
		}

		sensorsInDistrict, err := s.getDistrictSensorsFromCache(ctx, district.Id)
		if err != nil || sensorsInDistrict == nil {
			zap.L().Error("failed to get sensors by districtId", zap.Error(err), zap.Int64("districtId", district.Id))
			continue
		}
		wg.Go(func() { getTrustedSensor(ctx, respChan, *sensorsInDistrict) })
	}
	wg.Wait()
	close(respChan)

	var sensors []model.Sensor
	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	select {
	case s.cSensors <- sensors:
	case <-ctx.Done():
	}
}

func getTrustedSensor(ctx context.Context, resChan chan model.Sensor, sensors []sensor.Sensor) {
	var syncSensorList model.SyncSensorsList
	syncSensorList.Wg.Add(len(sensors))
	for _, sensor := range sensors {
		if ctx.Err() != nil {
			syncSensorList.Wg.Done()
			continue
		}
		go request.GetArchiveSensor(ctx, &syncSensorList, sensor.ApiId, sensor.District.Name)
	}
	syncSensorList.Wg.Wait()

	trustedAqlSensor := syncSensorList.GetSensor()
	if trustedAqlSensor == nil {
		return
	}

	select {
	case resChan <- *trustedAqlSensor:
	case <-ctx.Done():
	}
}
