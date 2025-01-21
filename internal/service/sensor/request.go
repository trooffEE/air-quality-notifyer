package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"air-quality-notifyer/internal/lib"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	endpoint = "https://airkemerovo.ru/api/sensor/current/%d?client_secret=guest"
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

func (s *SyncAirqualitySensorList) findWorstSensor() AqiSensor {
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
		getLastUpdatedSensor(&syncSensorList, sensor.ApiId, sensor.District.Name)
	}
	syncSensorList.wg.Wait()

	worstAirqualitySensor := syncSensorList.findWorstSensor()

	resChan <- worstAirqualitySensor
}

func fetchSensorById(id int64) (AqiSensorResponse, error) {
	res, err := http.Get(fmt.Sprintf(endpoint, id))
	if err != nil {
		lib.LogError("fetchSensorById", "failed to fetch sensor with id of %d", err, id)
		return AqiSensorResponse{}, nil
	}
	defer res.Body.Close()

	var aqiSensorsResponse AqiSensorResponse
	err = json.NewDecoder(res.Body).Decode(&aqiSensorsResponse)
	if err != nil {
		lib.LogError("fetchSensorById", "failed to decode response with status code %d", err, res.StatusCode)
		return AqiSensorResponse{}, nil
	}
	return aqiSensorsResponse, nil
}

func getLastUpdatedSensor(syncSensorList *SyncAirqualitySensorList, id int64, districtName string) {
	defer syncSensorList.wg.Done()

	response, err := fetchSensorById(id)
	if err != nil {
		lib.LogError("getLastUpdatedSensor", "failed to fetch sensor with id of %d", err, id)
		return
	}
	archivedSensors := response.Archive

	if len(archivedSensors) > 0 {
		latestDataFromSensor := archivedSensors[0]

		latestDataFromSensor.withDistrict(districtName)
		latestDataFromSensor.withApiData(id)

		syncSensorList.addSensor(latestDataFromSensor)
	}
}
