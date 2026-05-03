package commander

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/constants"
	"air-quality-notifyer/internal/helper"
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	CallbackTextFaqFromSetup          = "Узнать подробнее"
	CallbackDataFaqFromSetup          = "operation_mode_faq__operating-mode"
	CallbackTextSetup                 = "🌿 Режимы работы"
	CallbackDataSetup                 = "operation_mode"
	CallbackTextSetCity               = "Город 🏙"
	CallbackDataSetCity               = "set_operation_mode_city"
	CallbackTextAskForDistrictOptions = "Район 🏘"
	CallbackDataAskForDistrictOptions = "set_operation_mode_district"
	CallbackTextSetHome               = "Дом 🏡"
	CallbackDataSetHome               = "set_operation_mode_home"
	CallbackTextBack                  = "Назад"
	CallbackTextBackToSetup           = "🌿 Вернуться к настройкам режимов работы"
)

func NewModeCallbackHandlersRegistry(c *Commander) HandlersRegistry {
	return HandlersRegistry{
		CallbackDataFaqFromSetup:          c.Faq,
		CallbackDataSetup:                 c.Setup,
		CallbackDataSetCity:               c.SetCity,
		CallbackDataAskForDistrictOptions: c.AskForDistrictOptions,
	}
}

func (c *Commander) Setup(ctx context.Context, update tgbotapi.Update) {
	fmt.Println("Setup")
	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("Пожалуйста, выберите один из трех режимов работы для его настройки:\n\nЕсли не знайте какой режим выбрать, нажмите на \"%s\", чтобы получить информацию о них", CallbackTextFaqFromSetup),
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CallbackTextSetCity, CallbackDataSetCity),
			tgbotapi.NewInlineKeyboardButtonData(CallbackTextAskForDistrictOptions, CallbackDataAskForDistrictOptions),
			c.API.HomeMapInlineKeyboardButton(CallbackTextSetHome),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CallbackTextFaqFromSetup, CallbackDataFaqFromSetup),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CallbackTextBack, CallbackDataBack),
		),
	)

	if err := c.API.Edit(ctx, api.EditMessageConfig{Msg: msg, Markup: &markup}); err != nil {
		zap.L().Error("Error sending operating_mode message", zap.Error(err))
	}
}

func (c *Commander) Faq(ctx context.Context, update tgbotapi.Update) {
	markup := tgbotapi.NewInlineKeyboardMarkup()

	if update.CallbackQuery.Data == CallbackDataFaqFromSetup {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CallbackTextBackToSetup, CallbackDataSetup),
		))
	}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(CallbackTextBack, CallbackDataBack),
	))

	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf(
			"⚙️ <strong>Режимы работы</strong> ⚙️\n\n"+
				"🏙 <i>Город</i> 🏙\n\n"+
				"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков по <strong>всему городу</strong>. Данный функционал следует использовать, если вы хотите следить за общим состоянием воздуха в городе 🍃\n\n\n"+
				"🏘 <i>Район</i> 🏘\n\n"+
				"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков по <strong>выбранному району</strong> Кемерово. Данный функционал следует использовать, если вы хотите следить за конекретным районом/районами города 🍃\n\n\n"+
				"🏡 <i>Дом</i> 🏡\n\n"+
				"Режим работы, в котором бот отслеживает и отправляет оповещения от датчиков <strong>в пределах километра от выбранного места на карте или выбранных в ручную вами</strong>. Данный функционал следует использовать, если вы хотите следить за конкретными интересующими датчиками 🍃\n\n",
		),
	)

	if err := c.API.Edit(ctx, api.EditMessageConfig{Msg: msg, Markup: &markup}); err != nil {
		zap.L().Error("Error sending operating_mode_faq message", zap.Error(err))
	}
}

func (c *Commander) SetCity(ctx context.Context, update tgbotapi.Update) {
	// City mode is implemented end-to-end: mode is saved and user gets a confirmation.
	message := update.CallbackQuery.Message
	chatId := message.Chat.ID
	err := c.Services.User.SetOperatingMode(ctx, chatId, constants.City)
	if err != nil {
		zap.L().Error("Error setting operating mode", zap.Error(err))
		return
	}

	msg := tgbotapi.NewMessage(
		chatId,
		"🏙 Город 🏙\n\nТеперь вы будете получать оповещения с датчиков по всему городу! 🍃",
	)

	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending Mode.Set message", zap.Error(err))
	}

	if err = c.API.Delete(ctx, message); err != nil {
		zap.L().Error("Error deleting prev message", zap.Error(err))
	}
}

func (c *Commander) AskForDistrictOptions(ctx context.Context, update tgbotapi.Update) {
	// District mode setup currently stops at poll creation.
	// TODO: Persist selected districts to users_observed_districts and set constants.District as operating mode.
	chatID := update.CallbackQuery.Message.Chat.ID
	districts := c.Services.District.GetAllDistrictsNames(ctx)
	response, err := c.API.SendPoll(ctx, chatID, api.PollConfig{
		Question: "Интересующие районы 🏘:",
		Options:  districts,
	})
	if err != nil {
		zap.L().Error("set district: error sending poll", zap.Error(err))
		return
	}
	c.Services.District.SaveDistrictPollMessageInCache(
		ctx,
		response.Poll.ID,
		response.Chat.ID,
		response.MessageID,
	)
}

func (c *Commander) HandleDistrictsOptionsResult(ctx context.Context, pollUpdate *tgbotapi.Poll) {
	if pollUpdate == nil || pollUpdate.TotalVoterCount == 0 {
		return
	}

	cachedPollState, err := c.Services.District.GetDistrictPollMessageInCache(ctx, pollUpdate.ID)
	if err != nil {
		zap.L().Error("failed to get poll state from cache", zap.Error(err), zap.String("pollId", pollUpdate.ID))
		return
	}

	if err = c.API.DeleteTrackedMessageByOffset(ctx, cachedPollState.ChatID, 1); err != nil {
		zap.L().Error("Error deleting district setup message", zap.Error(err))
	}

	messageToDelete := tgbotapi.NewDeleteMessage(cachedPollState.ChatID, cachedPollState.MessageID)
	if err = c.API.DeleteRequest(ctx, messageToDelete); err != nil {
		zap.L().Error("Error sending DeleteMessage", zap.Error(err))
	}

	selectedOptions := helper.Filter(pollUpdate.Options, func(item tgbotapi.PollOption) bool { return item.VoterCount > 0 })
	if len(selectedOptions) == 0 {
		msg := tgbotapi.NewMessage(cachedPollState.ChatID, "Для режима \"Район\" нужно выбрать хотя бы один район. Попробуйте снова в настройках.")
		if sendErr := c.API.Send(ctx, api.MessageConfig{Msg: msg}); sendErr != nil {
			zap.L().Error("failed to send district mode empty selection message", zap.Error(sendErr))
		}
		return
	}

	allDistricts := c.Services.District.GetAllDistricts(ctx)
	districtIDByName := make(map[string]int64, len(allDistricts))
	for _, district := range allDistricts {
		districtIDByName[district.Name] = district.Id
	}

	selectedDistrictIDs := make([]int64, 0, len(selectedOptions))
	selectedDistrictNames := make([]string, 0, len(selectedOptions))
	for _, option := range selectedOptions {
		districtID, exists := districtIDByName[option.Text]
		if !exists {
			zap.L().Warn("poll option does not match district", zap.String("districtName", option.Text), zap.String("pollId", pollUpdate.ID))
			continue
		}

		selectedDistrictIDs = append(selectedDistrictIDs, districtID)
		selectedDistrictNames = append(selectedDistrictNames, option.Text)
	}

	if len(selectedDistrictIDs) == 0 {
		zap.L().Error("failed to map district poll options to district IDs", zap.String("pollId", pollUpdate.ID))
		return
	}

	err = c.Services.User.SetObservedDistricts(ctx, cachedPollState.ChatID, selectedDistrictIDs)
	if err != nil {
		zap.L().Error("failed to save observed districts", zap.Error(err), zap.Int64("chatId", cachedPollState.ChatID))
		return
	}

	err = c.Services.User.SetOperatingMode(ctx, cachedPollState.ChatID, constants.District)
	if err != nil {
		zap.L().Error("failed to set district operating mode", zap.Error(err), zap.Int64("chatId", cachedPollState.ChatID))
		return
	}

	msg := tgbotapi.NewMessage(
		cachedPollState.ChatID,
		fmt.Sprintf(
			"🏘 Район 🏘\n\nТеперь вы будете получать оповещения по выбранным районам! 🍃\n\nВы выбрали:\n%s",
			strings.Join(selectedDistrictNames, "\n"),
		),
	)
	if sendErr := c.API.Send(ctx, api.MessageConfig{Msg: msg}); sendErr != nil {
		zap.L().Error("failed to send district mode confirmation", zap.Error(sendErr))
	}
}
