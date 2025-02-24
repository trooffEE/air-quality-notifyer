package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"math"
	"net/http"
	"slices"
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

func (s *SyncAirqualitySensorList) sortAqi() {
	slices.SortFunc(s.list, func(a, b AqiSensor) int {
		if a.Aqi < b.Aqi {
			return -1
		} else if a.Aqi > b.Aqi {
			return 1
		}
		return 0
	})
}

// "Trusted" - meaning median from what we have in district, slightly more realistic than AVG and Worst AQI value determination
func (s *SyncAirqualitySensorList) getTrustedAqiSensor() *AqiSensor {
	if len(s.list) == 0 {
		return nil
	}
	s.sortAqi()
	trustedIndex := math.Ceil(float64(len(s.list) / 2))
	return &s.list[int(trustedIndex)]
}

func findTrustedSensor(resChan chan AqiSensor, sensors []models.AirqualitySensor) {
	var syncSensorList SyncAirqualitySensorList
	syncSensorList.wg.Add(len(sensors))
	for _, sensor := range sensors {
		go getLastUpdatedSensor(&syncSensorList, sensor.ApiId, sensor.District.Name)
	}
	syncSensorList.wg.Wait()

	trustedAqlSensor := syncSensorList.getTrustedAqiSensor()
	if trustedAqlSensor == nil {
		return
	}

	resChan <- *trustedAqlSensor
}

func fetchSensorById(id int64) (AqiSensorResponse, error) {
	res, err := http.Get(fmt.Sprintf(endpoint, id))
	if err != nil {
		zap.L().Error("failed to fetch sensor", zap.Error(err), zap.Int64("sensorId", id))
		return AqiSensorResponse{}, nil
	}
	defer res.Body.Close()

	var aqiSensorsResponse AqiSensorResponse
	err = json.NewDecoder(res.Body).Decode(&aqiSensorsResponse)
	if err != nil {
		zap.L().Error("failed to decode response with status code", zap.Error(err), zap.Int("statusCode", res.StatusCode))
		return AqiSensorResponse{}, nil
	}
	return aqiSensorsResponse, nil
}

func getLastUpdatedSensor(syncSensorList *SyncAirqualitySensorList, id int64, districtName string) {
	defer syncSensorList.wg.Done()

	response, err := fetchSensorById(id)
	if err != nil {
		zap.L().Error("failed to fetch sensor by id", zap.Error(err), zap.Int64("sensorId", id))
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
