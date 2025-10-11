package request

import (
	"air-quality-notifyer/internal/service/sensor/model"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var (
	endpoint = "https://airkemerovo.ru/api/sensor/current/%d?client_secret=guest"
)

type Response struct {
	Id          int            `json:"id"`
	CityId      int            `json:"city_id"`
	Description interface{}    `json:"description"`
	Lat         float64        `json:"lat"`
	Lon         float64        `json:"lon"`
	Address     string         `json:"address"`
	Floor       int            `json:"floor"`
	Radius      int            `json:"radius"`
	Source      interface{}    `json:"source"`
	Type        string         `json:"type"`
	Last        model.Sensor   `json:"last"`
	Archive     []model.Sensor `json:"archive"`
}

func GetArchiveSensor(syncSensors *model.SyncSensorsList, id int64, districtName string) {
	defer syncSensors.Wg.Done()

	response, err := fetchSensorById(id)
	if err != nil {
		zap.L().Error("failed to fetch sensor by id", zap.Error(err), zap.Int64("sensorId", id))
		return
	}
	sensors := response.Archive

	if len(sensors) > 0 {
		latestDataFromSensor := sensors[0]

		latestDataFromSensor.WithDistrict(districtName)
		latestDataFromSensor.WithApiData(id)

		syncSensors.AddSensor(latestDataFromSensor)
	}
}

func fetchSensorById(id int64) (Response, error) {
	res, err := http.Get(fmt.Sprintf(endpoint, id))
	if err != nil {
		zap.L().Error("failed to fetch sensor", zap.Error(err), zap.Int64("sensorId", id))
		return Response{}, nil
	}
	defer res.Body.Close()

	var sensorResponse Response
	if err = json.NewDecoder(res.Body).Decode(&sensorResponse); err != nil {
		zap.L().Error("failed to decode response with status code", zap.Error(err), zap.Int("statusCode", res.StatusCode))
		return Response{}, nil
	}
	return sensorResponse, nil
}
