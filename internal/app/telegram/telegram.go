package telegram

import (
	"air-quality-notifyer/internal/app/telegram/commander"
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/constants"
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
	modeHandlers := map[constants.ModeType]func([]model.Sensor){
		constants.City:     t.notifyCityUsers,
		constants.District: t.notifyDistrictUsers,
		constants.Home:     t.notifyHomeUsers,
	}

	modeOrder := []constants.ModeType{constants.City, constants.District, constants.Home}
	for _, mode := range modeOrder {
		modeHandlers[mode](sensors)
	}
}

func (t *tgBot) notifyCityUsers(sensors []model.Sensor) {
	cityUsers := t.Commander.Services.User.GetUsersIdsByOperatingMode(constants.City)
	messages := newUserMessages(sensors)
	for _, userID := range cityUsers {
		t.sendMessagesToUser(userID, messages)
	}
}

func (t *tgBot) notifyDistrictUsers(sensors []model.Sensor) {
	userDistricts := t.Commander.Services.User.GetObservedDistrictIdsByOperatingMode(constants.District)
	if len(userDistricts) == 0 {
		return
	}

	districtNameByID := map[int64]string{}
	for _, district := range t.Commander.Services.District.GetAllDistricts() {
		districtNameByID[district.Id] = district.Name
	}

	for userID, observedDistrictIDs := range userDistricts {
		observedDistrictNames := map[string]struct{}{}
		for _, districtID := range observedDistrictIDs {
			districtName, exists := districtNameByID[districtID]
			if !exists {
				continue
			}
			observedDistrictNames[districtName] = struct{}{}
		}

		if len(observedDistrictNames) == 0 {
			continue
		}

		districtSensors := []model.Sensor{}
		for _, sensor := range sensors {
			if _, exists := observedDistrictNames[sensor.District]; exists {
				districtSensors = append(districtSensors, sensor)
			}
		}

		t.sendMessagesToUser(userID, newUserMessages(districtSensors))
	}
}

func (t *tgBot) notifyHomeUsers(_ []model.Sensor) {
	// Placeholder for future Home mode implementation.
}

func (t *tgBot) sendMessagesToUser(userID int64, messages []string) {
	for _, message := range messages {
		payload := api.MessageConfig{Msg: tgbotapi.NewMessage(userID, message)}
		if err := t.Commander.API.Send(payload); err != nil && err.Code == 403 {
			t.Commander.Services.User.Delete(userID)
			break
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
