package helper

import "air-quality-notifyer/internal/constants"

func IsValidMode(mode constants.ModeType) bool {
	return mode >= constants.City && mode <= constants.Home
}

func Filter[T any](items []T, fn func(item T) bool) []T {
	filteredItems := []T{}
	for _, value := range items {
		if fn(value) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}
