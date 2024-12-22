package districts

import (
	repo "air-quality-notifyer/internal/db/repository"
)

type Service struct {
	repo repo.DistrictRepositoryType
}

func NewDistrictService(ur repo.DistrictRepositoryType) *Service {
	return &Service{
		repo: ur,
	}
}

func (s *Service) GetDistrictByCoords(x, y float64) int64 {
	id := s.repo.GetAssociatedDistrictIdByCoords(x, y)
	return id
}
