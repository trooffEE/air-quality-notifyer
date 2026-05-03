package api

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const pendingCommandTTL = 24 * time.Hour

type command string

func (a *Api) HandlePendingCommand(ctx context.Context, chatID int64) string {
	if commandName, exists := a.GetPendingCommand(ctx, chatID); exists {
		a.DeletePendingCommand(ctx, chatID)
		return commandName
	}

	return ""
}

func (a *Api) SetPendingCommand(ctx context.Context, chatID int64, command string) {
	if a.cache == nil {
		zap.L().Error("pending command cache is not configured")
		return
	}

	err := a.cache.Set(ctx, pendingCommandKey(chatID), command, pendingCommandTTL).Err()
	if err != nil {
		zap.L().Error("failed to set pending command", zap.Error(err), zap.Int64("chatId", chatID), zap.String("command", command))
	}
}

func (a *Api) GetPendingCommand(ctx context.Context, chatID int64) (string, bool) {
	if a.cache == nil {
		zap.L().Error("pending command cache is not configured")
		return "", false
	}

	command, err := a.cache.Get(ctx, pendingCommandKey(chatID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", false
	}
	if err != nil {
		zap.L().Error("failed to get pending command", zap.Error(err), zap.Int64("chatId", chatID))
		return "", false
	}

	return command, command != ""
}

func (a *Api) DeletePendingCommand(ctx context.Context, chatID int64) {
	if a.cache == nil {
		zap.L().Error("pending command cache is not configured")
		return
	}

	err := a.cache.Del(ctx, pendingCommandKey(chatID)).Err()
	if err != nil {
		zap.L().Error("failed to delete pending command", zap.Error(err), zap.Int64("chatId", chatID))
	}
}

func pendingCommandKey(chatID int64) string {
	return "telegram:" + strconv.FormatInt(chatID, 10) + ":pending_command"
}
