package districts

import (
	"air-quality-notifyer/internal/db/repository/districts"
	"air-quality-notifyer/internal/db/repository/sensor"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Service struct {
	repo  districts.Interface
	cache *redis.Client
}

type Interface interface {
	GetAllDistricts() []districts.District
	GetAllDistrictsNames() []string
	GetDistrictByCoords(x, y float64) *sensor.DistrictSensor
	GetDistrictPollMessageInCache(pollID string) (*DistrictPollMessage, error)
	SaveDistrictPollMessageInCache(pollID string, chatID int64, messageId int)
}

func New(ur districts.Interface, cache *redis.Client) Interface {
	return &Service{
		repo:  ur,
		cache: cache,
	}
}

// TODO Think about seperated module for polls (at least not in one service like now)
type DistrictPollMessage struct {
	ChatID    int64  "json:\"chat_id\""
	MessageID int    "json:\"message_id\"" // Needed for post cleanup logic
	PollID    string "json:\"poll_id\""
}

func (s *Service) SaveDistrictPollMessageInCache(pollID string, chatID int64, messageId int) {
	key := DistrictPollCacheKey(pollID)
	value := DistrictPollMessage{
		MessageID: messageId,
		ChatID:    chatID,
		PollID:    pollID,
	}
	payload, err := json.Marshal(value)
	if err != nil {
		zap.L().Error("failed to marshal sensor", zap.Error(err), zap.Any("payload", payload))
		return
	}
	err = s.cache.Set(context.Background(), key, payload, time.Minute*30).Err()
	if err != nil {
		zap.L().Error("cache: failed to save districts options", zap.Error(err))
	}
}

func (s *Service) GetDistrictPollMessageInCache(pollID string) (*DistrictPollMessage, error) {
	key := DistrictPollCacheKey(pollID)
	result, err := s.cache.Get(context.Background(), key).Result()
	if err != nil {
		zap.L().Error("cache: failed to fetch districts options", zap.Error(err))
		return nil, err
	}

	var options DistrictPollMessage
	err = json.Unmarshal([]byte(result), &options)
	if err != nil {
		zap.L().Error("failed to unmarshal districts options", zap.Error(err))
		return nil, err
	}

	return &options, nil
}

func DistrictPollCacheKey(pollID string) string {
	return fmt.Sprintf("telegram:poll:%s:flag", pollID)
}

func (s *Service) GetDistrictByCoords(x, y float64) *sensor.DistrictSensor {
	return s.repo.GetAssociatedDistrictIdByCoords(x, y)
}

func (s *Service) GetAllDistricts() []districts.District {
	districtsList, err := s.repo.GetAllDistricts()
	if err != nil {
		zap.L().Panic("Failed to get all districts", zap.Error(err))
	}

	return districtsList
}

func (s *Service) GetAllDistrictsNames() []string {
	districtsList, err := s.repo.GetAllDistrictsNames()
	if err != nil {
		zap.L().Panic("Failed to get all districts", zap.Error(err))
	}

	return districtsList
}
