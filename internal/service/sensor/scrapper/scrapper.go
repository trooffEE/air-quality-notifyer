package scrapper

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type Sensor struct {
	Id        int64   `json:"sensor_id"`
	Address   string  `json:"address"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	CreatedAt string  `json:"created_at"`
}

var setLastSensorsDataScriptStringStart, setLastSensorsDataScriptStringEnd = "setLastData('", "');"

func Scrap(ctx context.Context) ([]Sensor, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://airkemerovo.ru", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("access airkemerovo.ru: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		zap.L().Info("failed to grasp new sensors, airkemerovo page responded with ", zap.Int("status", res.StatusCode))
		return []Sensor{}, nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("parse airkemerovo.ru: %w", err)
	}

	var scriptContents string
	doc.Find("script[type='application/javascript']").Each(func(i int, s *goquery.Selection) {
		sText := s.Text()
		if strings.Contains(sText, setLastSensorsDataScriptStringStart) {
			scriptContents = sText
		}
	})
	reader := strings.NewReader(strings.TrimSpace(scriptContents))
	scanner := bufio.NewScanner(reader)

	var sensors []Sensor
	for scanner.Scan() {
		scriptLine := scanner.Text()

		if strings.Index(scriptLine, setLastSensorsDataScriptStringStart) != -1 {
			startJsonIndex := strings.Index(scriptLine, setLastSensorsDataScriptStringStart) + len(setLastSensorsDataScriptStringStart)
			endJsonIndex := strings.Index(scriptLine, setLastSensorsDataScriptStringEnd)

			jsonString := scriptLine[startJsonIndex:endJsonIndex]

			if err := json.Unmarshal([]byte(jsonString), &sensors); err != nil {
				zap.L().Error("failed to unmarshal json string", zap.Error(err))
			}
		}
	}

	return sensors, nil
}

// FilterSensorsByHourDiff TODO Maybe not a place for it
func FilterSensorsByHourDiff(sensors []Sensor, allowedDiffInHours int) []Sensor {
	var aliveSensors []Sensor
	layout := "2006-01-02T15:04:05.999999999Z"
	for _, sensor := range sensors {
		sensorTime, err := time.Parse(layout, sensor.CreatedAt)

		if err != nil {
			zap.L().Error("failed to parse time from sensor.CreatedAt", zap.Error(err), zap.Any("sensor", sensor))
		}

		diffInHours := sensorTime.Sub(time.Now().UTC()).Hours()
		if -1*diffInHours <= float64(allowedDiffInHours) {
			aliveSensors = append(aliveSensors, sensor)
		}
	}

	return aliveSensors
}
