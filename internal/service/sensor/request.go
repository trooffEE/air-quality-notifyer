package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"air-quality-notifyer/internal/lib"
	"encoding/json"
	"fmt"
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
	defer res.Body.Close()
	if err != nil {
		lib.LogError("fetchSensorById", "failed to fetch sensor with id of %d", err, id)
		return
	}

	var aqiSensorsResponse AqiSensorResponse
	err = json.NewDecoder(res.Body).Decode(&aqiSensorsResponse)
	if err != nil {
		lib.LogError("fetchSensorById", "failed to decode response with status code %d", err, res.StatusCode)
		return
	}

	archivedSensors := aqiSensorsResponse.Archive

	if len(archivedSensors) > 0 {
		latestDataFromSensor := archivedSensors[0]

		latestDataFromSensor.withDistrict(districtName)
		latestDataFromSensor.withApiData(id)

		syncSensorList.addSensor(latestDataFromSensor)
	}
}
