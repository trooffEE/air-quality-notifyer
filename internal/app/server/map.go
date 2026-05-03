package server

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"air-quality-notifyer/internal/config"
	"air-quality-notifyer/internal/constants"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

const telegramInitDataHeader = "X-Telegram-Init-Data"

type mapHandler struct {
	cfg      config.Config
	services Services
}

type homeSelectionRequest struct {
	SensorAPIIDs []int64 `json:"sensor_api_ids"`
}

type telegramInitUser struct {
	ID int64 `json:"id"`
}

func newMapHandler(cfg config.Config, services Services) *mapHandler {
	return &mapHandler{
		cfg:      cfg,
		services: services,
	}
}

func (h *mapHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(api.HomeMapPath, h.handleHomeMap)
	mux.HandleFunc(api.AliveSensorsPath, h.handleAliveSensors)
	mux.HandleFunc(api.HomeMapSelectionPath, h.handleHomeSelection)
}

func (h *mapHandler) handleHomeMap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	http.ServeFile(w, r, "frontend/home-map.html")
}

func (h *mapHandler) handleAliveSensors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if _, err := h.telegramUserIDFromRequest(r); err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sensors, err := h.services.Sensor.GetAliveSensorsFromCache(r.Context())
	if err != nil {
		zap.L().Error("failed to get alive sensors for map", zap.Error(err))
		http.Error(w, "failed to load sensors", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, sensors)
}

func (h *mapHandler) handleHomeSelection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	telegramUserID, err := h.telegramUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var payload homeSelectionRequest
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if len(payload.SensorAPIIDs) == 0 {
		http.Error(w, "sensor_api_ids must not be empty", http.StatusBadRequest)
		return
	}

	if err = h.validateAliveSensorSelection(r.Context(), payload.SensorAPIIDs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.services.User.SetObservedSensorsByAPIIds(r.Context(), telegramUserID, payload.SensorAPIIDs); err != nil {
		zap.L().Error("failed to save home sensors", zap.Error(err), zap.Int64("telegramUserID", telegramUserID))
		http.Error(w, "failed to save sensors", http.StatusInternalServerError)
		return
	}

	if err = h.services.User.SetOperatingMode(r.Context(), telegramUserID, constants.Home); err != nil {
		zap.L().Error("failed to set home operating mode", zap.Error(err), zap.Int64("telegramUserID", telegramUserID))
		http.Error(w, "failed to set home mode", http.StatusInternalServerError)
		return
	}

	h.sendHomeConfirmation(r.Context(), telegramUserID, len(payload.SensorAPIIDs))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *mapHandler) validateAliveSensorSelection(ctx context.Context, sensorAPIIDs []int64) error {
	aliveSensors, err := h.services.Sensor.GetAliveSensorsFromCache(ctx)
	if err != nil {
		zap.L().Error("failed to validate home sensor selection", zap.Error(err))
		return errors.New("failed to validate sensors")
	}

	aliveSensorIDs := make(map[int64]struct{}, len(aliveSensors))
	for _, sensor := range aliveSensors {
		aliveSensorIDs[sensor.APIID] = struct{}{}
	}

	for _, sensorAPIID := range sensorAPIIDs {
		if sensorAPIID <= 0 {
			return errors.New("sensor_api_ids must contain positive ids")
		}
		if _, exists := aliveSensorIDs[sensorAPIID]; !exists {
			return errors.New("sensor_api_ids must contain only active sensors")
		}
	}

	return nil
}

func (h *mapHandler) telegramUserIDFromRequest(r *http.Request) (int64, error) {
	initData := r.Header.Get(telegramInitDataHeader)
	if initData == "" {
		return 0, errors.New("missing telegram init data")
	}

	valid, err := tgbotapi.ValidateWebAppData(h.cfg.App.TelegramToken, initData)
	if err != nil || !valid {
		return 0, errors.New("invalid telegram init data")
	}

	return telegramUserIDFromInitData(initData)
}

func telegramUserIDFromInitData(initData string) (int64, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return 0, err
	}

	rawUser := values.Get("user")
	if rawUser == "" {
		return 0, errors.New("missing telegram user")
	}

	var user telegramInitUser
	if err = json.Unmarshal([]byte(rawUser), &user); err != nil {
		return 0, err
	}
	if user.ID == 0 {
		return 0, errors.New("missing telegram user id")
	}

	return user.ID, nil
}

func (h *mapHandler) sendHomeConfirmation(ctx context.Context, chatID int64, sensorsCount int) {
	if h.services.Bot == nil {
		return
	}

	msg := tgbotapi.NewMessage(
		chatID,
		"🏡 Дом 🏡\n\nТеперь вы будете получать оповещения по выбранным датчикам! 🍃",
	)
	if sensorsCount == 1 {
		msg.Text = "🏡 Дом 🏡\n\nТеперь вы будете получать оповещения по выбранному датчику! 🍃"
	}
	msg.ParseMode = tgbotapi.ModeHTML

	if err := h.services.Bot.Send(ctx, api.MessageConfig{Msg: msg}); err != nil {
		zap.L().Error("failed to send home mode confirmation", zap.Any("error", err), zap.Int64("chatID", chatID))
		return
	}

	if err := h.services.Bot.DeleteTrackedMessageByOffset(ctx, chatID, 1); err != nil {
		zap.L().Error("failed to delete home setup message", zap.Error(err), zap.Int64("chatID", chatID))
	}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		zap.L().Error("failed to write json response", zap.Error(err))
	}
}
