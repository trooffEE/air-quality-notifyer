package bot

import (
	"fmt"
	"strings"
)

type mentionSlug = string

type mentionResponse struct {
	Response string
}

const (
	notCommandMessage mentionSlug = "notCommandMessage"
)

var mentions = map[mentionSlug]mentionResponse{
	notCommandMessage: {
		Response: fmt.Sprintf("😓Пожалуйста, на данный момент я понимаю только команды, начинающиейся на символ \"/\":\n %s", strings.Join(publicCommandsList, "\n")),
	},
}

func getMessageByMention(mention mentionSlug) string {
	if extractedMention, ok := mentions[mention]; ok {
		return extractedMention.Response
	}
	return ""
}

func getMessageWithAQIStatsForChosenDistrict() string {
	// AQI = ((AQI_high - AQI_low) / (Conc_high - Conc_low)) * (Conc_measured - Conc_low) + AQI_low
	return ""
}
