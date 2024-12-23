package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type SyncAirqualitySensorList struct {
	mu   sync.Mutex
	wg   sync.WaitGroup
	list []AirqualitySensor
}

func (s *SyncAirqualitySensorList) addSensor(sensor AirqualitySensor) {
	s.mu.Lock()
	s.list = append(s.list, sensor)
	s.mu.Unlock()
}

func (s *SyncAirqualitySensorList) findWorstAirqualitySensor() AirqualitySensor {
	var worstAQISensor AirqualitySensor
	var currentWorstAQI float64
	for _, value := range s.list {
		if currentWorstAQI < value.AQIPM10 || currentWorstAQI < value.AQIPM25 {
			worstAQISensor = value
		}

		if currentWorstAQI < value.AQIPM25 {
			currentWorstAQI = value.AQIPM25
		} else if currentWorstAQI < value.AQIPM10 {
			currentWorstAQI = value.AQIPM10
		}
	}

	return worstAQISensor
}

func findWorstSensorInDistrict(resChan chan AirqualitySensor, sensors []models.AirqualitySensor) {
	var syncSensorList SyncAirqualitySensorList
	syncSensorList.wg.Add(len(sensors))

	for _, sensor := range sensors {
		fetchSensorById(&syncSensorList, sensor.ApiId, sensor.District.Name)
	}
	syncSensorList.wg.Wait()

	worstAirqualitySensor := syncSensorList.findWorstAirqualitySensor()

	resChan <- worstAirqualitySensor
}

func fetchSensorById(syncSensorList *SyncAirqualitySensorList, id int64, districtName string) {
	defer syncSensorList.wg.Done()

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

	var fetchedSensorsList []AirqualitySensor
	err = json.NewDecoder(res.Body).Decode(&fetchedSensorsList)
	if err != nil {
		log.Println("Something went wrong on decoding JSON from API step")
	}

	if len(fetchedSensorsList) > 0 {
		latestUpdatedSensor := fetchedSensorsList[len(fetchedSensorsList)-1]

		latestUpdatedSensor.withApiData(id)
		latestUpdatedSensor.withDistrict(districtName)

		syncSensorList.addSensor(latestUpdatedSensor)
	}
}
