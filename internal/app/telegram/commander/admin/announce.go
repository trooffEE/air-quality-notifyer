package admin

import (
	"air-quality-notifyer/internal/app/telegram/commander/api"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"go.uber.org/zap"
)

const (
	CommandAnnounce    = "/announce"
	announcementHeader = "🤖\n\n"
)

func (c *Commander) Announce(update tgbotapi.Update) {
	if !c.api.IsAdmin(update) {
		return
	}

	text, entities := announcementPayload(update.Message)
	if text == "" {
		return
	}
	text, entities = announcementMessage(text, entities)

	for _, userID := range c.service.User.GetUsersIds() {
		msg := tgbotapi.NewMessage(userID, text)
		msg.Entities = entities

		if err := c.api.Send(api.MessageConfig{Msg: msg}); err != nil {
			if err.Code == 403 {
				c.service.User.Delete(userID)
				continue
			}

			zap.L().Error("Error sending announcement", zap.Error(err), zap.Int64("userId", userID))
		}
	}
}

func IsAnnounceCommand(text string) bool {
	_, ok := announceCommandEnd(text)
	return ok
}

func announcementPayload(message *tgbotapi.Message) (string, []tgbotapi.MessageEntity) {
	if message == nil {
		return "", nil
	}

	commandEnd, ok := announceCommandEnd(message.Text)
	if !ok {
		return "", nil
	}

	payloadStart := len(message.Text)
	for offset, r := range message.Text[commandEnd:] {
		if !unicode.IsSpace(r) {
			payloadStart = commandEnd + offset
			break
		}
	}

	text := message.Text[payloadStart:]
	if text == "" {
		return "", nil
	}

	return text, shiftEntities(message.Entities, utf16Len(message.Text[:payloadStart]))
}

func announcementMessage(text string, entities []tgbotapi.MessageEntity) (string, []tgbotapi.MessageEntity) {
	return announcementHeader + text, addEntityOffset(entities, utf16Len(announcementHeader))
}

func announceCommandEnd(text string) (int, bool) {
	if !strings.HasPrefix(text, CommandAnnounce) {
		return 0, false
	}

	end := len(CommandAnnounce)
	if end < len(text) && text[end] == '@' {
		botNameStart := end + 1
		end = botNameStart

		for end < len(text) {
			r, size := utf8.DecodeRuneInString(text[end:])
			if !isBotUsernameRune(r) {
				break
			}
			end += size
		}

		if end == botNameStart {
			return 0, false
		}
	}

	if end == len(text) {
		return end, true
	}

	r, _ := utf8.DecodeRuneInString(text[end:])
	return end, unicode.IsSpace(r)
}

func isBotUsernameRune(r rune) bool {
	return r == '_' || r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}

func shiftEntities(entities []tgbotapi.MessageEntity, offset int) []tgbotapi.MessageEntity {
	if len(entities) == 0 {
		return nil
	}

	shifted := make([]tgbotapi.MessageEntity, 0, len(entities))
	for _, entity := range entities {
		start := entity.Offset
		end := entity.Offset + entity.Length
		if end <= offset {
			continue
		}

		if start < offset {
			entity.Offset = 0
			entity.Length = end - offset
		} else {
			entity.Offset = start - offset
		}

		if entity.Length > 0 {
			shifted = append(shifted, entity)
		}
	}

	return shifted
}

func addEntityOffset(entities []tgbotapi.MessageEntity, offset int) []tgbotapi.MessageEntity {
	if len(entities) == 0 {
		return nil
	}

	shifted := make([]tgbotapi.MessageEntity, len(entities))
	copy(shifted, entities)

	for index := range shifted {
		shifted[index].Offset += offset
	}

	return shifted
}

func utf16Len(text string) int {
	return len(utf16.Encode([]rune(text)))
}
