package commander

import (
	"context"
	"fmt"
	"slices"

	"air-quality-notifyer/internal/app/telegram/commander/api"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	CallbackTextFaq  = "❓ Режимы работы"
	CallbackDataFaq  = "operation_mode_faq"
	CallbackDataBack = "back"
)

const (
	CommandSettings = "⚙️ Настройки"
	CommandFaq      = "❓ FAQ"
)

var options = []string{CommandFaq, CommandSettings, CommandShowUsers, CommandPing}

func NewMenuMessageHandlersRegistry(c *Commander) HandlersRegistry {
	return HandlersRegistry{
		CommandFaq:      c.MenuFaq,
		CommandSettings: c.Settings,
	}
}

func NewMenuCallbackHandlersRegistry(c *Commander) HandlersRegistry {
	return HandlersRegistry{
		CallbackDataBack: c.MenuBack,
		CallbackDataFaq:  c.Faq,
	}
}

func (c *Commander) MenuBack(ctx context.Context, update tgbotapi.Update) {
	if err := c.API.Delete(ctx, update.CallbackQuery.Message); err != nil {
		zap.L().Error("Error deleting previous menu message", zap.Error(err))
	}
}

func (c *Commander) MenuFaq(ctx context.Context, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		"⚙️<strong>Ответы на вопросы</strong>\n\n"+
			"<i>- Связан ли данный бот с https://airkemerovo.ru ?</i>\n"+
			"Данный проект берет информацию именно с этого сервиса, используя публичный API, однако разработчик бота <strong>никак не связан</strong> с https://airkemerovo.ru\n\n"+
			"<i>- Бесплатно ли использование?</i>\n"+
			"Да, проект доступен для использования всем желающим, кто хочет следить за состоянием воздуху и получать уведомления об этом 🌿\n\n"+
			"<i>- Что умеет данный бот?</i>\n"+
			"Оповещать об опасном воздухе в разных режимах\n"+
			"1. по всему городу (режим - \"город\")\n"+
			"2. * по районам города (режим - \"район\")\n"+
			"3. * по датчикам (режим - \"дом\")\n\n"+
			"* - функционал в работе\n"+
			"ℹ Подробнее о режимах можно прочитать в меню - \n\"%s\"\n\n"+
			"<i>- Почему на https://airkemerovo.ru больше датчиков, чем в предложенном перечне?</i>\n"+
			"Потому что на данном этапе проекта реализована работа только с <strong>датчиками города Кемерово</strong>. Это сознательное решение автора проекта, однако это вполне может поменяться в будущем\n\n",
		CallbackTextFaq,
	),
	)
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CallbackTextFaq, CallbackDataFaq),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", CallbackDataBack),
		),
	)

	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg, Markup: markup}); err != nil {
		zap.L().Error("Error sending faq message", zap.Error(err))
	}
}

func IsMenuButton(button string) bool {
	return slices.Contains(options, button)
}

func (c *Commander) Settings(ctx context.Context, update tgbotapi.Update) {
	if ctx.Err() != nil {
		return
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"⚙️ <strong>Настройки</strong>\n"+
			"Здесь вы можете настроить нужный функционал бота",
	)

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CallbackTextSetup, CallbackDataSetup),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(CallbackTextBack, CallbackDataBack),
		),
	)

	if err := c.API.Send(ctx, api.MessageConfig{Msg: msg, Markup: markup}); err != nil {
		zap.L().Error("Error sending configure message", zap.Error(err))
	}
}
