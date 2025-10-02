package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var (
	endpoint = "https://airkemerovo.ru/api/sensor/current/%d?client_secret=guest"
)

func getLastUpdatedSensor(syncSensors *SyncTrustedSensors, id int64, districtName string) {
	defer syncSensors.wg.Done()

	response, err := fetchSensorById(id)
	if err != nil {
		zap.L().Error("failed to fetch sensor by id", zap.Error(err), zap.Int64("sensorId", id))
		return
	}
	sensors := response.Archive

	if len(sensors) > 0 {
		latestDataFromSensor := sensors[0]

		latestDataFromSensor.withDistrict(districtName)
		latestDataFromSensor.withApiData(id)

		syncSensors.addSensor(latestDataFromSensor)
	}
}

func fetchSensorById(id int64) (SensorResponse, error) {
	res, err := http.Get(fmt.Sprintf(endpoint, id))
	if err != nil {
		zap.L().Error("failed to fetch sensor", zap.Error(err), zap.Int64("sensorId", id))
		return SensorResponse{}, nil
	}
	defer res.Body.Close()

	var sensorResponse SensorResponse
	if err = json.NewDecoder(res.Body).Decode(&sensorResponse); err != nil {
		zap.L().Error("failed to decode response with status code", zap.Error(err), zap.Int("statusCode", res.StatusCode))
		return SensorResponse{}, nil
	}
	return sensorResponse, nil
}
