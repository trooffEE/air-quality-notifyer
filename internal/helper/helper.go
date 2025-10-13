package helper

import "air-quality-notifyer/internal/constants"

func IsValidMode(mode constants.ModeType) bool {
	return mode >= constants.City && mode <= constants.Home
}
