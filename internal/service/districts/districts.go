package districts

import (
	"air-quality-notifyer/internal/db/models"
	repo "air-quality-notifyer/internal/db/repository"

	"go.uber.org/zap"
)

type Service struct {
	repo repo.DistrictRepositoryInterface
}

func New(ur repo.DistrictRepositoryInterface) *Service {
	return &Service{
		repo: ur,
	}
}

func (s *Service) GetDistrictByCoords(x, y float64) *models.DistrictSensor {
	return s.repo.GetAssociatedDistrictIdByCoords(x, y)
}

func (s *Service) GetAllDistricts() []models.District {
	districtsList, err := s.repo.GetAllDistricts()
	if err != nil {
		zap.L().Panic("Failed to get all districts", zap.Error(err))
	}

	return districtsList
}
