package api

import (
	"fmt"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

const (
	HomeMapPath          = "/miniapp/home-map"
	AliveSensorsPath     = "/api/map/alive-sensors"
	HomeMapSelectionPath = "/api/map/home-selection"
)

func (a *Api) HomeMapURL() string {
	return strings.TrimRight(a.miniAppUrl(), "/") + HomeMapPath
}

func (a *Api) miniAppUrl() string {
	if a.cfg.App.MiniAppUrl != "" {
		return a.cfg.App.MiniAppUrl
	}

	return fmt.Sprintf("http://localhost:%s", a.cfg.App.HttpServerPort)
}

func (a *Api) HomeMapWebAppInfo() tgbotapi.WebAppInfo {
	return tgbotapi.WebAppInfo{URL: a.HomeMapURL()}
}

func (a *Api) HomeMapInlineKeyboardButton(text string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonWebApp(text, a.HomeMapWebAppInfo())
}
