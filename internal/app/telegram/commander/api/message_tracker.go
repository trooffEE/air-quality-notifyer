package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	trackedMessageTTL   = 5 * time.Minute
	trackedMessageLimit = 20
)

func trackMessage(ctx context.Context, cache *redis.Client, chatID int64, messageID int) {
	if cache == nil || messageID == 0 {
		return
	}

	key := trackedMessagesCacheKey(chatID)
	messageIDValue := strconv.Itoa(messageID)

	pipe := cache.TxPipeline()
	pipe.LRem(ctx, key, 0, messageIDValue)
	pipe.LPush(ctx, key, messageIDValue)
	pipe.LTrim(ctx, key, 0, trackedMessageLimit-1)
	pipe.Expire(ctx, key, trackedMessageTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		zap.L().Error("cache: failed to track telegram message", zap.Error(err), zap.Int64("chatId", chatID), zap.Int("messageId", messageID))
	}
}

func deleteTrackedMessageByOffset(ctx context.Context, cache *redis.Client, bot *tgbotapi.BotAPI, chatID int64, offset int64) error {
	if cache == nil || bot == nil {
		return nil
	}

	messageIDValue, err := cache.LIndex(ctx, trackedMessagesCacheKey(chatID), offset).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("get tracked telegram message: %w", err)
	}

	messageID, err := strconv.Atoi(messageIDValue)
	if err != nil {
		return fmt.Errorf("parse tracked telegram message id %q: %w", messageIDValue, err)
	}

	if _, err = bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID)); err != nil {
		return fmt.Errorf("delete tracked telegram message %d: %w", messageID, err)
	}

	untrackMessage(ctx, cache, chatID, messageID)
	return nil
}

func untrackMessage(ctx context.Context, cache *redis.Client, chatID int64, messageID int) {
	if cache == nil || messageID == 0 {
		return
	}

	if err := cache.LRem(ctx, trackedMessagesCacheKey(chatID), 0, strconv.Itoa(messageID)).Err(); err != nil {
		zap.L().Error("cache: failed to untrack telegram message", zap.Error(err), zap.Int64("chatId", chatID), zap.Int("messageId", messageID))
	}
}

func (a *Api) DeleteTrackedMessages(ctx context.Context, chatID int64, count int) {
	for range count {
		if err := a.DeleteTrackedMessageByOffset(ctx, chatID, 0); err != nil {
			zap.L().Error("failed to delete tracked telegram message", zap.Error(err), zap.Int64("chatId", chatID))
			return
		}
	}
}

func (a *Api) DeleteTrackedMessageByOffset(ctx context.Context, chatID int64, offset int64) error {
	return deleteTrackedMessageByOffset(ctx, a.cache, a.Bot, chatID, offset)
}

func (a *Api) trackMessage(ctx context.Context, chatID int64, messageID int) {
	trackMessage(ctx, a.cache, chatID, messageID)
}

func (a *Api) untrackMessage(ctx context.Context, chatID int64, messageID int) {
	untrackMessage(ctx, a.cache, chatID, messageID)
}

func trackedMessagesCacheKey(chatID int64) string {
	return fmt.Sprintf("telegram:chat:%d:messages", chatID)
}
