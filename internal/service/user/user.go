package user

import (
	"air-quality-notifyer/internal/db/models"
	repo "air-quality-notifyer/internal/db/repository"
	"errors"

	"go.uber.org/zap"
)

type Service struct {
	repo repo.UserRepositoryInterface
}

func New(ur repo.UserRepositoryInterface) *Service {
	return &Service{
		repo: ur,
	}
}

func (ur *Service) IsNewUser(id int64) bool {
	_, err := ur.repo.FindById(id)

	if err != nil {
		if errors.Is(repo.UserNotFound, err) {
			return true
		}
		zap.L().Error("repository error", zap.Error(err))
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
		zap.L().Error("failed to register new user", zap.Error(err))
	}
}

func (ur *Service) GetUsersIds() []int64 {
	ids, err := ur.repo.GetAllIds()

	if err != nil {
		zap.L().Error("failed to get users ids", zap.Error(err))
	}

	return ids
}

func (ur *Service) GetUsersNames() []string {
	names, err := ur.repo.GetAllNames()

	if err != nil {
		zap.L().Error("failed to get users names", zap.Error(err))
	}

	return names
}

func (ur *Service) DeleteUser(id int64) {
	err := ur.repo.DeleteUserById(id)

	if err != nil {
		zap.L().Error("failed to delete user", zap.Error(err), zap.Int64("userId", id))
	}
}
