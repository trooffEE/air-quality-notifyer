package sensor

import (
	repo "air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/service/districts"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo       repo.SensorRepositoryInterface
	sDistricts *districts.Service
	cSensors   chan []AqiSensor
	syncCron   chan interface{}
	cache      *redis.Client
}

func NewSensorService(
	repo repo.SensorRepositoryInterface,
	sDistricts *districts.Service,
	cache *redis.Client,
) *Service {
	return &Service{
		repo:       repo,
		sDistricts: sDistricts,
		cSensors:   make(chan []AqiSensor),
		syncCron:   make(chan interface{}),
		cache:      cache,
	}
}

func (s *Service) ListenChangesInSensors(handler func([]AqiSensor)) {
	for update := range s.cSensors {
		handler(update)
	}
}
