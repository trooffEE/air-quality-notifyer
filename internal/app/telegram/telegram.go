package telegram

import (
	"air-quality-notifyer/internal/app/telegram/commander"
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/service/sensor/model"

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
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	return &tgBot{
		bot:       bot,
		updates:   bot.GetUpdatesChan(updateConfig),
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
