package districts

import (
	repo "air-quality-notifyer/internal/db/repository"
	"fmt"
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
	id, err := s.repo.GetAssociatedDistrictIdByCoords(x, y)
	if err != nil {
		fmt.Println(err)
	}
	return id
}
