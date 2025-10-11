package districts

import (
	"air-quality-notifyer/internal/db/repository/districts"
	"air-quality-notifyer/internal/db/repository/sensor"

	"go.uber.org/zap"
)

type Service struct {
	repo districts.Interface
}

func New(ur districts.Interface) *Service {
	return &Service{
		repo: ur,
	}
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
