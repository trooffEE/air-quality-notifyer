package message

import (
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

func IsCommand(message *tgbotapi.Message, command string) bool {
	if message == nil {
		return false
	}

	if message.Command() == normalizeCommand(command) {
		return true
	}

	return IsCommandText(message.Text, command)
}

func IsCommandText(text string, command string) bool {
	_, ok := commandEnd(text, command)
	return ok
}

func CommandPayload(message *tgbotapi.Message, command string) (string, []tgbotapi.MessageEntity, bool) {
	if message == nil {
		return "", nil, false
	}

	commandEnd, ok := commandEnd(message.Text, command)
	if !ok {
		return "", nil, false
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
		return "", nil, true
	}

	return text, ShiftEntities(message.Entities, UTF16Len(message.Text[:payloadStart])), true
}

func Prepend(prefix string, text string, entities []tgbotapi.MessageEntity) (string, []tgbotapi.MessageEntity) {
	return prefix + text, AddEntityOffset(entities, UTF16Len(prefix))
}

func ShiftEntities(entities []tgbotapi.MessageEntity, offset int) []tgbotapi.MessageEntity {
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

func AddEntityOffset(entities []tgbotapi.MessageEntity, offset int) []tgbotapi.MessageEntity {
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

func UTF16Len(text string) int {
	return len(utf16.Encode([]rune(text)))
}

func commandEnd(text string, command string) (int, bool) {
	prefix := "/" + normalizeCommand(command)
	if !strings.HasPrefix(text, prefix) {
		return 0, false
	}

	end := len(prefix)
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

func normalizeCommand(command string) string {
	return strings.TrimPrefix(command, "/")
}

func isBotUsernameRune(r rune) bool {
	return r == '_' || r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}
