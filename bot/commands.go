package bot

type command = string

const (
	help                     command = "/help"
	showAQIForChosenDistrict command = "/showAQI"
)

var publicCommandsDictionary map[command]command = map[command]command{
	help:                     help,
	showAQIForChosenDistrict: showAQIForChosenDistrict,
}

var publicCommandsList = []string{help, showAQIForChosenDistrict}

func isPublicCommandProvided(message string) bool {
	_, ok := publicCommandsDictionary[message]
	return ok
}

func isShowAQIForChosenDistrictCommandProvided(command string) bool {
	return command == showAQIForChosenDistrict
}

// Список комманд полезных для разработки, доступны исключительно разработчику по ID
var developmentCommand = []string{""}
