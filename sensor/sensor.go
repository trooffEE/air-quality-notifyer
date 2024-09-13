package sensor

import (
	"air-quality-notifyer/districts"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"sync"
)

// TODO вынести в отдельный пакет
type requestDone struct {
	_sensors []Data
}

func (r requestDone) notifyChangesInSensors() {
	NotifyChangesInSensors(r._sensors)
}

func GetSensorsDataOnceIn(cronString string) {
	sensors := NewSensorsData()
	c := cron.New()
	_, err := c.AddFunc(cronString, func() {
		fetchSensors(sensors).notifyChangesInSensors()
	})
	if err != nil {
		log.Panic(err)
	}
	c.Start()
}

func fetchSensors(sensors []Data) requestDone {
	respChan := make(chan Data, len(districts.Dictionary))

	for _, district := range districts.Dictionary {
		fetchSensorById(respChan, district)
	}

	close(respChan)

	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	return requestDone{sensors}
}

func richSensorData(
	fetchedSensorData Data,
	districtName string,
	id int,
) Data {
	sensorData := fetchedSensorData

	sensorData.District = districtName
	sensorData.Id = id
	sensorData.SourceLink = fmt.Sprintf("https://airkemerovo.ru/sensor/%d", id)
	if sensorData.Humidity >= 90 {
		sensorData.AdditionalInfo = "Высокая влажность. Показания PM могут быть не корректны\n"
	}
	if sensorData.Temperature < -60 {
		sensorData.AdditionalInfo += fmt.Sprintf("Датчики температуры в районе %s не исправен!\n", sensorData.District)
	}
	sensorData.getInformationAboutAQI()

	return sensorData
}

func fetchSensorById(resChan chan Data, district districts.DictionaryWithSensors) {
	var wg sync.WaitGroup
	wg.Add(len(district.SensorIds))

	var result []Data

	for _, id := range district.SensorIds {
		go func() {
			defer wg.Done()

			var fetchedSensorData []Data
			res, err := http.Post(
				fmt.Sprintf("https://airkemerovo.ru/api/sensor/archive/%d/1", id),
				"application/json",
				nil,
			)
			if err != nil {
				log.Printf("Error in API call for sensor ID %d: %v", id, err)
				return
			}
			defer res.Body.Close()

			if res == nil {
				log.Printf("Unexpected nil response from API for sensor ID: %d", id)
				return
			}

			err = json.NewDecoder(res.Body).Decode(&fetchedSensorData)
			if err != nil {
				log.Println("Something went wrong on decoding JSON from API step")
			}

			if len(fetchedSensorData) > 0 {
				result = append(result, richSensorData(fetchedSensorData[len(fetchedSensorData)-1], district.Name, id))
			}
		}()
	}

	wg.Wait()

	var worstAQISensor Data
	var currentWorstAQI float64
	for _, value := range result {
		if currentWorstAQI < value.AQIPM10 || currentWorstAQI < value.AQIPM25 {
			worstAQISensor = value
		}

		if currentWorstAQI < value.AQIPM25 {
			currentWorstAQI = value.AQIPM25
		} else if currentWorstAQI < value.AQIPM10 {
			currentWorstAQI = value.AQIPM10
		}
	}

	resChan <- worstAQISensor
}
