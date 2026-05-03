package telegram

import (
	"air-quality-notifyer/internal/app/telegram/commander"
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/constants"
	"air-quality-notifyer/internal/service/sensor/model"
	"context"
	"sync"

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
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	return &tgBot{
		bot:       bot,
		updates:   bot.GetUpdatesChan(updateConfig),
		Commander: cmder,
	}
}

func (t *tgBot) Start(ctx context.Context) func(context.Context) {
	var wg sync.WaitGroup
	wg.Go(func() { t.listenUpdates(ctx) })
	wg.Go(func() { t.listenSensors(ctx) })

	return func(shutdownCtx context.Context) {
		t.bot.StopReceivingUpdates()

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-shutdownCtx.Done():
			zap.L().Warn("timed out waiting for telegram listeners to stop", zap.Error(shutdownCtx.Err()))
		}
	}
}

func (t *tgBot) BotAPI() *tgbotapi.BotAPI {
	return t.bot
}

func (t *tgBot) listenSensors(ctx context.Context) {
	t.Commander.Services.Sensor.ListenChanges(ctx, t.notifyUsers)
}

func (t *tgBot) notifyUsers(ctx context.Context, sensors []model.Sensor) {
	modeHandlers := map[constants.ModeType]func(context.Context, []model.Sensor){
		constants.City:     t.notifyCityUsers,
		constants.District: t.notifyDistrictUsers,
		constants.Home:     t.notifyHomeUsers,
	}

	modeOrder := []constants.ModeType{constants.City, constants.District, constants.Home}
	for _, mode := range modeOrder {
		if ctx.Err() != nil {
			return
		}
		modeHandlers[mode](ctx, sensors)
	}
}

func (t *tgBot) notifyCityUsers(ctx context.Context, sensors []model.Sensor) {
	cityUsers := t.Commander.Services.User.GetUsersIdsByOperatingMode(ctx, constants.City)
	messages := newUserMessages(sensors)
	for _, userID := range cityUsers {
		t.sendMessagesToUser(ctx, userID, messages)
	}
}

func (t *tgBot) notifyDistrictUsers(ctx context.Context, sensors []model.Sensor) {
	userDistricts := t.Commander.Services.User.GetObservedDistrictIdsByOperatingMode(ctx, constants.District)
	if len(userDistricts) == 0 {
		return
	}

	districtNameByID := map[int64]string{}
	for _, district := range t.Commander.Services.District.GetAllDistricts(ctx) {
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

		t.sendMessagesToUser(ctx, userID, newUserMessages(districtSensors))
	}
}

func (t *tgBot) notifyHomeUsers(ctx context.Context, sensors []model.Sensor) {
	userSensors := t.Commander.Services.User.GetObservedSensorAPIIdsByOperatingMode(ctx, constants.Home)
	if len(userSensors) == 0 {
		return
	}

	for userID, observedSensorIDs := range userSensors {
		observedSensors := map[int64]struct{}{}
		for _, sensorID := range observedSensorIDs {
			observedSensors[sensorID] = struct{}{}
		}

		homeSensors := []model.Sensor{}
		for _, sensor := range sensors {
			if _, exists := observedSensors[sensor.Id]; exists {
				homeSensors = append(homeSensors, sensor)
			}
		}

		t.sendMessagesToUser(ctx, userID, newUserMessages(homeSensors))
	}
}

func (t *tgBot) sendMessagesToUser(ctx context.Context, userID int64, messages []string) {
	for _, message := range messages {
		if ctx.Err() != nil {
			return
		}

		payload := api.MessageConfig{Msg: tgbotapi.NewMessage(userID, message)}
		if err := t.Commander.API.Send(ctx, payload); err != nil && err.Code == 403 {
			t.Commander.Services.User.Delete(ctx, userID)
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

func (t *tgBot) listenUpdates(ctx context.Context) {
	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "🌀 Перезапустить бота",
		},
		tgbotapi.BotCommand{
			Command:     commander.CommandFeedback,
			Description: "💬 Обратная связь",
		},
	)

	if _, err := t.bot.Request(cfg); err != nil {
		zap.L().Error("commander request error", zap.Error(err))
	}

	t.Commander.HandleUpdate(ctx, t.updates)
}
