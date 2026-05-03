package user

import (
	"air-quality-notifyer/internal/constants"
	"air-quality-notifyer/internal/db/repository/user"
	"air-quality-notifyer/internal/helper"
	"air-quality-notifyer/internal/service/user/model"
	"context"
	"errors"

	"go.uber.org/zap"
)

type Service struct {
	repo user.Interface
}

type Interface interface {
	IsNew(ctx context.Context, id int64) bool
	Delete(ctx context.Context, id int64)
	GetUsersNames(ctx context.Context) []string
	GetUsersIds(ctx context.Context) []int64
	GetUsersIdsByOperatingMode(ctx context.Context, mode constants.ModeType) []int64
	GetObservedDistrictIdsByOperatingMode(ctx context.Context, mode constants.ModeType) map[int64][]int64
	GetObservedSensorAPIIdsByOperatingMode(ctx context.Context, mode constants.ModeType) map[int64][]int64
	Register(ctx context.Context, userModel model.User)
	SetOperatingMode(ctx context.Context, id int64, mode constants.ModeType) error
	SetObservedDistricts(ctx context.Context, id int64, districtIDs []int64) error
	SetObservedSensorsByAPIIds(ctx context.Context, id int64, sensorAPIIDs []int64) error
}

func New(ur user.Interface) Interface {
	return &Service{
		repo: ur,
	}
}

func (ur *Service) IsNew(ctx context.Context, id int64) bool {
	_, err := ur.repo.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, user.NotFound) {
			return true
		}
		zap.L().Error("repository error", zap.Error(err))
	}

	return false
}

func (ur *Service) Register(ctx context.Context, userModel model.User) {
	dto := user.User{
		TelegramId: userModel.Id,
		Username:   userModel.Username,
	}

	err := ur.repo.Register(ctx, dto)

	if err != nil {
		zap.L().Error("failed to register new user", zap.Error(err))
	}
}

func (ur *Service) GetUsersIds(ctx context.Context) []int64 {
	ids, err := ur.repo.GetAllIds(ctx)

	if err != nil {
		zap.L().Error("failed to get users ids", zap.Error(err))
	}

	return ids
}

func (ur *Service) GetUsersIdsByOperatingMode(ctx context.Context, mode constants.ModeType) []int64 {
	ids, err := ur.repo.GetAllIdsByOperatingMode(ctx, mode)
	if err != nil {
		zap.L().Error("failed to get users ids by operating mode", zap.Error(err), zap.Int("mode", mode))
	}

	return ids
}

func (ur *Service) GetObservedDistrictIdsByOperatingMode(ctx context.Context, mode constants.ModeType) map[int64][]int64 {
	observedDistricts, err := ur.repo.GetObservedDistrictIdsByOperatingMode(ctx, mode)
	if err != nil {
		zap.L().Error("failed to get observed districts by operating mode", zap.Error(err), zap.Int("mode", mode))
		return map[int64][]int64{}
	}

	return observedDistricts
}

func (ur *Service) GetObservedSensorAPIIdsByOperatingMode(ctx context.Context, mode constants.ModeType) map[int64][]int64 {
	observedSensors, err := ur.repo.GetObservedSensorAPIIdsByOperatingMode(ctx, mode)
	if err != nil {
		zap.L().Error("failed to get observed sensors by operating mode", zap.Error(err), zap.Int("mode", mode))
		return map[int64][]int64{}
	}

	return observedSensors
}

func (ur *Service) GetUsersNames(ctx context.Context) []string {
	names, err := ur.repo.GetAllNames(ctx)

	if err != nil {
		zap.L().Error("failed to get users names", zap.Error(err))
	}

	return names
}

func (ur *Service) Delete(ctx context.Context, id int64) {
	err := ur.repo.DeleteUserById(ctx, id)

	if err != nil {
		zap.L().Error("failed to delete user", zap.Error(err), zap.Int64("userId", id))
	}
}

func (ur *Service) SetOperatingMode(ctx context.Context, id int64, mode constants.ModeType) error {
	if !helper.IsValidMode(mode) {
		err := errors.New("invalid operating mode")
		zap.L().Error("Setting mode", zap.Error(err))
		return err
	}

	err := ur.repo.SetOperatingMode(ctx, id, mode)
	if err != nil {
		zap.L().Error("failed to set operating mode", zap.Error(err))
		return err
	}

	return nil
}

func (ur *Service) SetObservedDistricts(ctx context.Context, id int64, districtIDs []int64) error {
	err := ur.repo.SetObservedDistricts(ctx, id, districtIDs)
	if err != nil {
		zap.L().Error("failed to set observed districts", zap.Error(err), zap.Int64("userId", id))
		return err
	}

	return nil
}

func (ur *Service) SetObservedSensorsByAPIIds(ctx context.Context, id int64, sensorAPIIDs []int64) error {
	err := ur.repo.SetObservedSensorsByAPIIds(ctx, id, sensorAPIIDs)
	if err != nil {
		zap.L().Error("failed to set observed sensors", zap.Error(err), zap.Int64("userId", id))
		return err
	}

	return nil
}
