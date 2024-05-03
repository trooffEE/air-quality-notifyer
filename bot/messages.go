package bot

import (
	"fmt"
	"strings"
)

type MentionSlug = string

type MentionResponse struct {
	Response string
}

const (
	NotCommandMessage MentionSlug = "notCommandMessage"
)

var Mentions = map[MentionSlug]MentionResponse{
	NotCommandMessage: {
		Response: fmt.Sprintf("😓Пожалуйста, на данный момент я понимаю только команды, начинающиейся на символ \"/\":\n %s", strings.Join(PublicCommandsList, "\n")),
	},
}

func GetMessageByMention(mention MentionSlug) string {
	if extractedMention, ok := Mentions[mention]; ok {
		return extractedMention.Response
	}
	return ""
}

func GetMessageWithAQIStatsForChosenDistrict() string {
	// AQI = ((AQI_high - AQI_low) / (Conc_high - Conc_low)) * (Conc_measured - Conc_low) + AQI_low
	return ""
}
