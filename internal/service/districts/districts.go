package districts

import repo "air-quality-notifyer/internal/db/repository"

type Service struct {
	repo repo.DistrictRepositoryType
}

func NewDistrictService(ur repo.DistrictRepositoryType) *Service {
	return &Service{
		repo: ur,
	}
}

func (s *Service) GetDistrictsPeriodically() {
	/** TODO uncomment me */
	//cronString := "0 0 * * *"
}
