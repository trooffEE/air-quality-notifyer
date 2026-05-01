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

func MergeMaps[M ~map[K]V, K comparable, V any](maps ...M) M {
	totalLen := 0
	for _, currentMap := range maps {
		totalLen += len(currentMap)
	}

	mergedMap := make(M, totalLen)
	for _, currentMap := range maps {
		for key, value := range currentMap {
			mergedMap[key] = value
		}
	}

	return mergedMap
}

func HasOverlappingKeys[M ~map[K]V, K comparable, V any](maps ...M) bool {
	seenKeys := make(map[K]struct{})
	for _, currentMap := range maps {
		for key := range currentMap {
			if _, exists := seenKeys[key]; exists {
				return true
			}
			seenKeys[key] = struct{}{}
		}
	}

	return false
}
