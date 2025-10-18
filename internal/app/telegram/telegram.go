package telegram

import (
	"air-quality-notifyer/internal/app/telegram/commander"
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/service/sensor/model"
	"fmt"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type tgBot struct {
	bot       *tgbotapi.BotAPI
	updates   tgbotapi.UpdatesChannel
	Commander *commander.Commander
}

func Init(cfg config.Config, services *commander.Services) *tgBot {
	bot, err := tgbotapi.NewBotAPI(cfg.App.TelegramToken)
	if err != nil {
		zap.L().Error("Filed to create new bot api", zap.Error(err))
		panic(err)
	}

	cmder := commander.New(cfg, bot, services)
	if cfg.Development {
		bot.Debug = true

		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30

		return &tgBot{
			bot:       bot,
			updates:   bot.GetUpdatesChan(updateConfig),
			Commander: cmder,
		}
	}

	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("https://%s/webhook%s", cfg.App.WebhookHost, bot.Token))
	if err != nil {
		zap.L().Panic("Filed to create new webhook", zap.Error(err))
	}
	_, err = bot.Request(wh)
	if err != nil {
		panic(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		zap.L().Panic("Failed to get webhook info", zap.Error(err))
	}

	if info.LastErrorDate != 0 {
		zap.L().Error("failed to init get info about webhook", zap.Error(err))
	}

	updates := bot.ListenForWebhook(fmt.Sprintf("/webhook%s", bot.Token))

	return &tgBot{
		bot:       bot,
		updates:   updates,
		Commander: cmder,
	}
}

func (t *tgBot) Start() {
	go t.listenUpdates()
	go t.listenSensors()
}

func (t *tgBot) listenSensors() {
	t.Commander.Services.Sensor.ListenChanges(t.notifyUsers)
}

func (t *tgBot) notifyUsers(sensors []model.Sensor) {
	messages := newUserMessages(sensors)

	ids := t.Commander.Services.User.GetUsersIds()
	for _, id := range ids {
		for _, message := range messages {
			payload := api.MessageConfig{Msg: tgbotapi.NewMessage(id, message)}
			if err := t.Commander.API.Send(payload); err != nil && err.Code == 403 {
				t.Commander.Services.User.Delete(id)
				break
			}
		}
	}
}

func newUserMessages(sensors []model.Sensor) []string {
	var messages []string
	for _, sensor := range sensors {
		if sensor.IsDangerousLevelDetected() {
			msg := sensor.DangerLevelText()
			messages = append(messages, msg)
		}
	}
	return messages
}

func (t *tgBot) listenUpdates() {
	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "🌀 Перезапустить бота",
		},
	)
	if _, err := t.bot.Request(cfg); err != nil {
		zap.L().Error("commander request error", zap.Error(err))
	}

	t.Commander.HandleUpdate(t.updates)
}
