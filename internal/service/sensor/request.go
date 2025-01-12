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
	list []AqiSensor
}

func (s *SyncAirqualitySensorList) addSensor(sensor AqiSensor) {
	s.mu.Lock()
	s.list = append(s.list, sensor)
	s.mu.Unlock()
}

func (s *SyncAirqualitySensorList) findWorstAirqualitySensor() AqiSensor {
	var worstAQISensor AqiSensor
	var currentWorstAQI int
	for _, value := range s.list {
		if currentWorstAQI < value.Aqi {
			currentWorstAQI = value.Aqi
			worstAQISensor = value
		}
	}

	return worstAQISensor
}

func findWorstSensorInDistrict(resChan chan AqiSensor, sensors []models.AirqualitySensor) {
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

	res, err := http.Get(fmt.Sprintf("https://airkemerovo.ru/api/sensor/current/%d?client_secret=guest", id))
	if err != nil {
		log.Printf("Error in API call for sensor ID %d: %+v", id, err)
		return
	}
	defer res.Body.Close()

	var aqiSensorsResponse AqiSensorResponse
	err = json.NewDecoder(res.Body).Decode(&aqiSensorsResponse)
	if err != nil {
		log.Println("Something went wrong on decoding JSON from API step")
	}

	archivedSensors := aqiSensorsResponse.Archive

	if len(archivedSensors) > 0 {
		latestDataFromSensor := archivedSensors[0]

		latestDataFromSensor.withDistrict(districtName)
		latestDataFromSensor.withApiData(id)

		syncSensorList.addSensor(latestDataFromSensor)
	}
}
