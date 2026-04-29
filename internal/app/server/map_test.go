package server

import (
	"air-quality-notifyer/internal/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramUserIDFromInitData(t *testing.T) {
	initData := url.Values{}
	initData.Set("user", `{"id":12345,"first_name":"Test"}`)

	userID, err := telegramUserIDFromInitData(initData.Encode())
	require.NoError(t, err)
	assert.Equal(t, int64(12345), userID)
}

func TestTelegramUserIDFromRequestValidatesInitData(t *testing.T) {
	const token = "secret-token"
	handler := newMapHandler(config.Config{
		App: config.AppConfig{TelegramToken: token},
	}, Services{})

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	req.Header.Set(telegramInitDataHeader, signedTelegramInitData(token, map[string]string{
		"auth_date": "1710000000",
		"user":      `{"id":456,"first_name":"Test"}`,
	}))

	userID, err := handler.telegramUserIDFromRequest(req)
	require.NoError(t, err)
	assert.Equal(t, int64(456), userID)
}

func TestTelegramUserIDFromRequestRejectsInvalidInitData(t *testing.T) {
	handler := newMapHandler(config.Config{
		App: config.AppConfig{TelegramToken: "secret-token"},
	}, Services{})

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	req.Header.Set(telegramInitDataHeader, "user=%7B%22id%22%3A456%7D&hash=bad")

	_, err = handler.telegramUserIDFromRequest(req)
	require.Error(t, err)
}

func signedTelegramInitData(token string, fields map[string]string) string {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	dataCheckString := make([]string, 0, len(fields))
	values := url.Values{}
	for _, key := range keys {
		dataCheckString = append(dataCheckString, key+"="+fields[key])
		values.Set(key, fields[key])
	}

	secret := hmac.New(sha256.New, []byte("WebAppData"))
	secret.Write([]byte(token))

	hash := hmac.New(sha256.New, secret.Sum(nil))
	hash.Write([]byte(strings.Join(dataCheckString, "\n")))
	values.Set("hash", hex.EncodeToString(hash.Sum(nil)))

	return values.Encode()
}
