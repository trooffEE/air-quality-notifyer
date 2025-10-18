package districts

import (
	"air-quality-notifyer/internal/db/repository/districts"
	"air-quality-notifyer/internal/db/repository/sensor"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Service struct {
	repo  districts.Interface
	cache *redis.Client
}

type Interface interface {
	GetOptionForDistrict(userId int)
	GetAllDistricts() []districts.District
	GetDistrictByCoords(x, y float64) *sensor.DistrictSensor
}

func New(ur districts.Interface, cache *redis.Client) Interface {
	return &Service{
		repo:  ur,
		cache: cache,
	}
}

func (s *Service) GetOptionForDistrict(userId int) {
	//s.cache.Get("setup-district:" + strconv.Itoa(userId))
}

func (s *Service) GetDistrictByCoords(x, y float64) *sensor.DistrictSensor {
	return s.repo.GetAssociatedDistrictIdByCoords(x, y)
}

func (s *Service) GetAllDistricts() []districts.District {
	districtsList, err := s.repo.GetAllDistricts()
	if err != nil {
		zap.L().Panic("Failed to get all districts", zap.Error(err))
	}

	return districtsList
}
