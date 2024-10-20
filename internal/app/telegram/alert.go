package telegram

import (
	"air-quality-notifyer/internal/config"
	"log"
	"strconv"
)

// TODO Создать нормальный модуль логирования
func (t *tgBot) AlertAdminWithPanic(exception interface{}) {

	id, err := strconv.Atoi(config.Cfg.AdminTelegramId)

	if err != nil {
		log.Println("[alert] I mean this should never happen, right?")
		return
	}

	switch exception.(type) {
	case error:
		t.Commander.DefaultSend(int64(id), exception.(error).Error())
	}

}
