package commander

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const pendingCommandTTL = 24 * time.Hour

func (c *Commander) SetPendingCommand(ctx context.Context, chatID int64, command string) {
	if c.Services.Cache == nil {
		zap.L().Error("pending command cache is not configured")
		return
	}

	err := c.Services.Cache.Set(ctx, pendingCommandKey(chatID), command, pendingCommandTTL).Err()
	if err != nil {
		zap.L().Error("failed to set pending command", zap.Error(err), zap.Int64("chatId", chatID), zap.String("command", command))
	}
}

func (c *Commander) PendingCommand(ctx context.Context, chatID int64) (string, bool) {
	if c.Services.Cache == nil {
		zap.L().Error("pending command cache is not configured")
		return "", false
	}

	command, err := c.Services.Cache.Get(ctx, pendingCommandKey(chatID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", false
	}
	if err != nil {
		zap.L().Error("failed to get pending command", zap.Error(err), zap.Int64("chatId", chatID))
		return "", false
	}

	return command, command != ""
}

func (c *Commander) DeletePendingCommand(ctx context.Context, chatID int64) {
	if c.Services.Cache == nil {
		zap.L().Error("pending command cache is not configured")
		return
	}

	err := c.Services.Cache.Del(ctx, pendingCommandKey(chatID)).Err()
	if err != nil {
		zap.L().Error("failed to delete pending command", zap.Error(err), zap.Int64("chatId", chatID))
	}
}

func pendingCommandKey(chatID int64) string {
	return "telegram:" + strconv.FormatInt(chatID, 10) + ":pending_command"
}
