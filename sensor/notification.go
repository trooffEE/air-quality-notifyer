package sensor

/* Warning: changesInAPIAppearedChannel should not be closed */
var changesInAPIAppearedChannel chan []Data = make(chan []Data)

func NotifyChangesInSensors(sensors []Data) {
	changesInAPIAppearedChannel <- sensors
}

func ListenChangesInSensors(handler func([]Data)) {
	for update := range changesInAPIAppearedChannel {
		handler(update)
	}
}
