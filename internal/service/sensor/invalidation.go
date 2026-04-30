package sensor

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"air-quality-notifyer/internal/service/sensor/scrapper"
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	InvalidationPeriod = 4
)

func (s *Service) StartInvalidatingSensorsPeriodically(ctx context.Context) func(context.Context) {
	cronCreator := cron.New()
	cronString := fmt.Sprintf("0 */%d * * *", InvalidationPeriod)

	_, err := cronCreator.AddFunc(cronString, func() {
		if err := s.startInvalidation(ctx, InvalidationPeriod); err != nil {
			zap.L().Error("failed to invalidate sensors", zap.Error(err))
			return
		}

		select {
		case s.syncCron <- struct{}{}:
		case <-ctx.Done():
		}
	})
	if err != nil {
		panic(err)
	}

	cronCreator.Start()

	return func(shutdownCtx context.Context) {
		stopCron(shutdownCtx, cronCreator)
	}
}

func (s *Service) startInvalidation(ctx context.Context, allowedHourDiff int) error {
	scrappedSensors, err := scrapper.Scrap(ctx)
	if err != nil {
		return err
	}

	aliveSensors := scrapper.FilterSensorsByHourDiff(scrappedSensors, allowedHourDiff)

	for _, scrappedSensor := range aliveSensors {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		s.saveSensor(ctx, scrappedSensor)
	}

	return nil
}

func (s *Service) saveSensor(ctx context.Context, scrappedSensor scrapper.Sensor) {
	district := s.sDistricts.GetDistrictByCoords(ctx, scrappedSensor.Lat, scrappedSensor.Lon)
	// TODO Не работаем с датчиками вне районов города
	if district == nil {
		return
	}

	payload := sensor.Sensor{
		DistrictId: district.Id,
		ApiId:      scrappedSensor.Id,
		Address:    scrappedSensor.Address,
		Lat:        scrappedSensor.Lat,
		Lon:        scrappedSensor.Lon,
		CreatedAt:  scrappedSensor.CreatedAt,
		District: sensor.DistrictSensor{
			Id:   district.Id,
			Name: district.Name,
		},
	}

	s.saveSensorInCache(ctx, payload)
}

func stopCron(ctx context.Context, cronCreator *cron.Cron) {
	stopped := cronCreator.Stop()
	select {
	case <-stopped.Done():
	case <-ctx.Done():
		zap.L().Warn("timed out waiting for cron jobs to stop", zap.Error(ctx.Err()))
	}
}
