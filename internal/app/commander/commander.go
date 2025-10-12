package commander

import (
	"air-quality-notifyer/internal/app/commander/admin"
	"air-quality-notifyer/internal/app/commander/api"
	"air-quality-notifyer/internal/app/commander/mode"
	"air-quality-notifyer/internal/config"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Commander struct {
	API   api.Interface
	Admin admin.Interface
	Mode  mode.Interface
}

func New(bot *tgbotapi.BotAPI, cfg config.Config) *Commander {
	apiCmder, err := api.NewApi(bot, cfg)
	if err != nil {
		zap.S().Fatalw("Failed to create api interface", "error", err)
		return nil
	}

	return &Commander{
		API:   apiCmder,
		Admin: admin.New(apiCmder),
		Mode:  mode.New(apiCmder),
	}
}
