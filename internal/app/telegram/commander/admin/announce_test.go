package admin

import (
	tgmessage "air-quality-notifyer/internal/app/telegram/commander/message"
	"testing"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/stretchr/testify/assert"
)

func TestIsAnnounceCommand(t *testing.T) {
	t.Parallel()

	assert.True(t, IsAnnounceCommand("/announce message"))
	assert.True(t, IsAnnounceCommand("/announce@air_quality_bot message"))
	assert.False(t, IsAnnounceCommand("/announcement message"))
}

func TestAnnouncementPayload(t *testing.T) {
	t.Parallel()

	text, entities := announcementPayload(&tgbotapi.Message{
		Text: "/announce formatted text",
		Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(CommandAnnounce)},
			{Type: "bold", Offset: 10, Length: 9},
		},
	})

	assert.Equal(t, "formatted text", text)
	assert.Equal(t, []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 0, Length: 9},
	}, entities)
}

func TestAnnouncementPayloadHandlesUTF16Offsets(t *testing.T) {
	t.Parallel()

	text, entities := announcementPayload(&tgbotapi.Message{
		Text: "/announce alert 🚨 now",
		Entities: []tgbotapi.MessageEntity{
			{Type: "italic", Offset: 19, Length: 3},
		},
	})

	assert.Equal(t, "alert 🚨 now", text)
	assert.Equal(t, []tgbotapi.MessageEntity{
		{Type: "italic", Offset: 9, Length: 3},
	}, entities)
}

func TestAnnouncementMessageAddsHeader(t *testing.T) {
	t.Parallel()

	text, entities := announcementMessage("formatted text", []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 0, Length: 9},
	})

	assert.Equal(t, "🤖\n\nformatted text", text)
	assert.Equal(t, []tgbotapi.MessageEntity{
		{Type: "bold", Offset: tgmessage.UTF16Len(announcementHeader), Length: 9},
	}, entities)
}
