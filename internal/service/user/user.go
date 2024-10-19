package user

import (
	"air-quality-notifyer/internal/db/exceptions"
	repo "air-quality-notifyer/internal/db/repository"
	"errors"
	"fmt"
)

type Service struct {
	repo repo.UserRepositoryInterface
}

func NewUserService(ur repo.UserRepositoryInterface) *Service {
	return &Service{
		repo: ur,
	}
}

func (ur *Service) IsNewUser(id int64) bool {
	_, err := ur.repo.FindById(id)
	if errors.Is(exceptions.UserNotFound, err) {
		return true
	}

	if err != nil {
		fmt.Println(err)
	}

	return false
}
