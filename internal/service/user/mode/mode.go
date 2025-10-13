package mode

import (
	"air-quality-notifyer/internal/constants"
	"air-quality-notifyer/internal/db/repository/user"
	"air-quality-notifyer/internal/helper"

	"go.uber.org/zap"
)

type Mode struct {
	repo user.Interface
}

type Interface interface {
	Set(mode int) error
}

func New(repo user.Interface) Interface {
	return &Mode{repo: repo}
}

func (m *Mode) Set(mode constants.ModeType) error {
	if !helper.IsValidMode(mode) {
		zap.L().Error("Setting mode", zap.Error(IsInvalid))
		return IsInvalid
	}

	m.repo.GetAllIds()

	return nil
}
