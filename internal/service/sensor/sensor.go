package sensor

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"air-quality-notifyer/internal/service/districts"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo       sensor.Interface
	sDistricts *districts.Service
	cSensors   chan []Sensor
	syncCron   chan interface{}
	cache      *redis.Client
}

func New(
	repo sensor.Interface,
	sDistricts *districts.Service,
	cache *redis.Client,
) *Service {
	return &Service{
		repo:       repo,
		sDistricts: sDistricts,
		cSensors:   make(chan []Sensor),
		syncCron:   make(chan interface{}),
		cache:      cache,
	}
}

func (s *Service) ListenChangesInSensors(handler func([]Sensor)) {
	for update := range s.cSensors {
		handler(update)
	}
}
