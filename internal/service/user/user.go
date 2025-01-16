package user

import (
	"air-quality-notifyer/internal/db/exceptions"
	"air-quality-notifyer/internal/db/models"
	repo "air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/lib"
	"errors"
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

	if err != nil {
		if errors.Is(exceptions.UserNotFound, err) {
			return true
		}
		lib.LogError("IsNewUser", "repository error", err)
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
		lib.LogError("Register", "failed to register new user", err)
	}
}

func (ur *Service) GetUsersIds() []int64 {
	ids, err := ur.repo.GetAllIds()

	if err != nil {
		lib.LogError("GetUsersIds", "failed to get users ids", err)
	}

	return ids
}

func (ur *Service) GetUsersNames() []string {
	names, err := ur.repo.GetAllNames()

	if err != nil {
		lib.LogError("GetUsersNames", "failed to get users names", err)
	}

	return names
}

func (ur *Service) DeleteUser(id int64) {
	err := ur.repo.DeleteUserById(id)

	if err != nil {
		lib.LogError("DeleteUser", "failed to delete user %d", err, id)
	}
}
