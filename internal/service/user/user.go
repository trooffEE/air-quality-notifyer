package user

import (
	"air-quality-notifyer/internal/constants"
	"air-quality-notifyer/internal/db/repository/user"
	"air-quality-notifyer/internal/exception"
	"air-quality-notifyer/internal/helper"
	"air-quality-notifyer/internal/service/user/model"
	"errors"

	"go.uber.org/zap"
)

type Service struct {
	repo user.Interface
}

type Interface interface {
	IsNewUser(id int64) bool
	DeleteUser(id int64)
	GetUsersNames() []string
	GetUsersIds() []int64
	Register(userModel model.User)
	SetOperatingMode(id int64, mode constants.ModeType) error
}

func New(ur user.Interface) Interface {
	return &Service{
		repo: ur,
	}
}

func (ur *Service) IsNewUser(id int64) bool {
	_, err := ur.repo.FindById(id)

	if err != nil {
		if errors.Is(user.NotFound, err) {
			return true
		}
		zap.L().Error("repository error", zap.Error(err))
	}

	return false
}

func (ur *Service) Register(userModel model.User) {
	dto := user.User{
		TelegramId: userModel.Id,
		Username:   userModel.Username,
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

func (ur *Service) SetOperatingMode(id int64, mode constants.ModeType) error {
	if !helper.IsValidMode(mode) {
		zap.L().Error("Setting mode", zap.Error(exception.InvalidOperatingMode))
		return exception.InvalidOperatingMode
	}

	err := ur.repo.SetOperatingMode(id, mode)
	if err != nil {
		zap.L().Error("failed to set operating mode", zap.Error(err))
		return err
	}

	return nil
}
