package user

import (
	"air-quality-notifyer/internal/db/exceptions"
	"air-quality-notifyer/internal/db/models"
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

func (ur *Service) Register(user User) {
	dto := models.User{
		TelegramId: user.Id,
		Username:   user.Username,
	}

	err := ur.repo.Register(dto)

	if err != nil {
		fmt.Println(err)
	}
}

func (ur *Service) GetUsersIds() *[]int64 {
	ids, err := ur.repo.GetAllIds()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ids)

	return ids
}
