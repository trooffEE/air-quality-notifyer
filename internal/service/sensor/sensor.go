package sensor

import (
	repo "air-quality-notifyer/internal/db/repository"
	"air-quality-notifyer/internal/districts"
	districts2 "air-quality-notifyer/internal/service/districts"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"sync"
)

type Service struct {
	sensorChannel chan []Data
	//TODO rename districts2 to be districts
	districts *districts2.Service
	repo      repo.SensorRepositoryType
}

func NewSensorService(repository repo.SensorRepositoryType, districtService *districts2.Service) *Service {
	return &Service{
		repo:          repository,
		districts:     districtService,
		sensorChannel: make(chan []Data),
	}
}

func (s *Service) ListenChangesInSensors(handler func([]Data)) {
	for update := range s.sensorChannel {
		handler(update)
	}
}

func (s *Service) ScrapSensorDataPeriodically() {
	sensors := NewSensorsData()
	c, cronString := cron.New(), "0 * * * *"
	_, err := c.AddFunc(cronString, func() {
		sensors := fetchSensors(sensors)
		s.sensorChannel <- sensors
	})
	if err != nil {
		log.Panic(err)
	}
	c.Start()
}

func fetchSensors(sensors []Data) []Data {
	respChan := make(chan Data, len(districts.Dictionary))

	for _, district := range districts.Dictionary {
		fetchSensorById(respChan, district)
	}

	close(respChan)

	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	return sensors
}

// TODO Needs refactor
func fetchSensorById(resChan chan Data, district districts.DictionaryWithSensors) {
	var wg sync.WaitGroup
	wg.Add(len(district.SensorIds))

	var result []Data
	var mu sync.Mutex

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
				log.Printf("Error in API call for sensor ID %d: %+v", id, err)
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
				mu.Lock()
				result = append(result, richSensorData(fetchedSensorData[len(fetchedSensorData)-1], district.Name, id))
				mu.Unlock()
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
