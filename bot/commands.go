package bot

type Command = string

const (
	help                     Command = "/help"
	showAQIForChosenDistrict Command = "/showAQI"
)

var PublicCommandsDictionary map[Command]Command = map[Command]Command{
	help:                     help,
	showAQIForChosenDistrict: showAQIForChosenDistrict,
}

var PublicCommandsList = []string{help, showAQIForChosenDistrict}

func IsPublicCommandProvided(message string) bool {
	_, ok := PublicCommandsDictionary[message]
	return ok
}

func isShowAQIForChosenDistrictCommandProvided(command string) bool {
	return command == showAQIForChosenDistrict
}

// Список комманд полезных для разработки, доступны исключительно разработчику по ID
var DevelopmentCommand = []string{""}
