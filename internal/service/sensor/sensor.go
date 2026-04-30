package sensor

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"air-quality-notifyer/internal/service/districts"
	"air-quality-notifyer/internal/service/sensor/model"
	"context"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo       sensor.Interface
	sDistricts districts.Interface
	cSensors   chan []model.Sensor
	syncCron   chan struct{}
	cache      *redis.Client
}

type AliveSensor struct {
	APIID    int64   `json:"api_id"`
	Address  string  `json:"address"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	District string  `json:"district"`
}

type Interface interface {
	ListenChanges(ctx context.Context, handler func(context.Context, []model.Sensor))
	StartInvalidatingSensorsPeriodically(ctx context.Context) func(context.Context)
	StartGettingTrustedSensorsEveryHour(ctx context.Context) func(context.Context)
	GetAliveSensorsFromCache(ctx context.Context) ([]AliveSensor, error)
}

func New(
	repo sensor.Interface,
	sDistricts districts.Interface,
	cache *redis.Client,
) Interface {
	return &Service{
		repo:       repo,
		sDistricts: sDistricts,
		cSensors:   make(chan []model.Sensor),
		syncCron:   make(chan struct{}),
		cache:      cache,
	}
}

func (s *Service) ListenChanges(ctx context.Context, handler func(context.Context, []model.Sensor)) {
	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-s.cSensors:
			if !ok {
				return
			}
			handler(ctx, update)
		}
	}
}
