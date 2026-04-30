package mode

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/constants"
	"air-quality-notifyer/internal/helper"
	sDistricts "air-quality-notifyer/internal/service/districts"
	sUser "air-quality-notifyer/internal/service/user"
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

type Commander struct {
	api     api.Interface
	service Service
}

type Service struct {
	User     sUser.Interface
	District sDistricts.Interface
}

type Interface interface {
	Setup(update tgbotapi.Update)
	Faq(update tgbotapi.Update)
	SetCity(ctx context.Context, update tgbotapi.Update)
	AskForDistrictOptions(ctx context.Context, update tgbotapi.Update)
	HandleDistrictsOptionsResult(ctx context.Context, update *tgbotapi.Poll)
}

func New(api api.Interface, service Service) Interface {
	return &Commander{
		api:     api,
		service: service,
	}
}

func (c *Commander) Setup(update tgbotapi.Update) {
	msg := tgbotapi.NewEditMessageText(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("Пожалуйста, выберите один из трех режимов работы для его настройки:\n\nЕсли не знайте какой режим выбрать, нажмите на \"%s\", чтобы получить информацию о них", KeypadFaqFromSetupText),
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(KeypadSetCityText, KeypadSetCityData),
			tgbotapi.NewInlineKeyboardButtonData(KeypadAskForDistrictOptionsText, KeypadAskForDistrictOptionsData),
			c.api.HomeMapInlineKeyboardButton(KeypadSetHomeText),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(KeypadFaqFromSetupText, KeypadFaqFromSetupData),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(api.KeypadMenuBackText, api.KeypadMenuBackData),
		),
	)

	if err := c.api.Edit(api.EditMessageConfig{Msg: msg, Markup: &markup}); err != nil {
		zap.L().Error("Error sending operating_mode message", zap.Error(err))
	}
}

func (c *Commander) Faq(update tgbotapi.Update) {
	markup := tgbotapi.NewInlineKeyboardMarkup()

	if update.CallbackQuery.Data == KeypadFaqFromSetupData {
		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(KeypadBackText, KeypadSetupData),
		))
	}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(api.KeypadMenuBackText, api.KeypadMenuBackData),
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

	if err := c.api.Edit(api.EditMessageConfig{Msg: msg, Markup: &markup}); err != nil {
		zap.L().Error("Error sending operating_mode_faq message", zap.Error(err))
	}
}

func (c *Commander) SetCity(ctx context.Context, update tgbotapi.Update) {
	// City mode is implemented end-to-end: mode is saved and user gets a confirmation.
	message := update.CallbackQuery.Message
	chatId := message.Chat.ID
	err := c.service.User.SetOperatingMode(ctx, chatId, constants.City)
	if err != nil {
		zap.L().Error("Error setting operating mode", zap.Error(err))
		return
	}

	msg := tgbotapi.NewMessage(
		chatId,
		"🏙 Город 🏙\n\nТеперь вы будете получать оповещения с датчиков по всему городу! 🍃",
	)

	if err := c.api.Send(api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("Error sending Mode.Set message", zap.Error(err))
	}

	if err = c.api.Delete(message); err != nil {
		zap.L().Error("Error deleting prev message", zap.Error(err))
	}
}

func (c *Commander) AskForDistrictOptions(ctx context.Context, update tgbotapi.Update) {
	// District mode setup currently stops at poll creation.
	// TODO: Persist selected districts to users_observed_districts and set constants.District as operating mode.
	chatID := update.CallbackQuery.Message.Chat.ID
	districts := c.service.District.GetAllDistrictsNames(ctx)
	response, err := c.api.SendPoll(chatID, api.PollConfig{
		Question: "Интересующие районы 🏘:",
		Options:  districts,
	})
	if err != nil {
		zap.L().Error("set district: error sending poll", zap.Error(err))
		return
	}
	c.service.District.SaveDistrictPollMessageInCache(
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

	cachedPollState, err := c.service.District.GetDistrictPollMessageInCache(ctx, pollUpdate.ID)
	if err != nil {
		zap.L().Error("failed to get poll state from cache", zap.Error(err), zap.String("pollId", pollUpdate.ID))
		return
	}

	if err = c.api.DeleteTrackedMessageByOffset(ctx, cachedPollState.ChatID, 1); err != nil {
		zap.L().Error("Error deleting district setup message", zap.Error(err))
	}

	messageToDelete := tgbotapi.NewDeleteMessage(cachedPollState.ChatID, cachedPollState.MessageID)
	if err = c.api.DeleteRequest(messageToDelete); err != nil {
		zap.L().Error("Error sending DeleteMessage", zap.Error(err))
	}

	selectedOptions := helper.Filter(pollUpdate.Options, func(item tgbotapi.PollOption) bool { return item.VoterCount > 0 })
	if len(selectedOptions) == 0 {
		msg := tgbotapi.NewMessage(cachedPollState.ChatID, "Для режима \"Район\" нужно выбрать хотя бы один район. Попробуйте снова в настройках.")
		if sendErr := c.api.Send(api.MessageConfig{Msg: msg}); sendErr != nil {
			zap.L().Error("failed to send district mode empty selection message", zap.Error(sendErr))
		}
		return
	}

	allDistricts := c.service.District.GetAllDistricts(ctx)
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

	err = c.service.User.SetObservedDistricts(ctx, cachedPollState.ChatID, selectedDistrictIDs)
	if err != nil {
		zap.L().Error("failed to save observed districts", zap.Error(err), zap.Int64("chatId", cachedPollState.ChatID))
		return
	}

	err = c.service.User.SetOperatingMode(ctx, cachedPollState.ChatID, constants.District)
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
	if sendErr := c.api.Send(api.MessageConfig{Msg: msg}); sendErr != nil {
		zap.L().Error("failed to send district mode confirmation", zap.Error(sendErr))
	}
}
