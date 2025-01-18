package sensor

import (
	"air-quality-notifyer/internal/lib"
	"bufio"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
)

type AqiSensorScriptScrapped struct {
	Id        int64   `json:"sensor_id"`
	Address   string  `json:"address"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	CreatedAt string  `json:"created_at"`
}

var setSensorsStringStart, setSensorsStringEnd = "setLastData('", "');"

func scrapSensorData() []AqiSensorScriptScrapped {
	res, err := http.Get("https://airkemerovo.ru")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		lib.LogMessage("scrapSensorData", "failed to grasp new sensors, airkemerovo page responded with %d", res.StatusCode)
		return []AqiSensorScriptScrapped{}
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var scriptContents string
	doc.Find("script[type='application/javascript']").Each(func(i int, s *goquery.Selection) {
		sText := s.Text()
		if strings.Contains(sText, setSensorsStringStart) {
			scriptContents = sText
		}
	})
	reader := strings.NewReader(strings.TrimSpace(scriptContents))
	scanner := bufio.NewScanner(reader)

	var sensors []AqiSensorScriptScrapped
	for scanner.Scan() {
		scriptLine := scanner.Text()

		if strings.Index(scriptLine, setSensorsStringStart) != -1 {
			startJsonIndex := strings.Index(scriptLine, setSensorsStringStart) + len(setSensorsStringStart)
			endJsonIndex := strings.Index(scriptLine, setSensorsStringEnd)

			jsonString := scriptLine[startJsonIndex:endJsonIndex]
			err := json.Unmarshal([]byte(jsonString), &sensors)
			if err != nil {
				lib.LogError("scrapSensorData", "failed to unmarshal json string ", err)
			}
		}
	}

	return sensors
}

func filterDeadSensors(sensors []AqiSensorScriptScrapped, allowedDiffInHours int) []AqiSensorScriptScrapped {
	aliveSensors := make([]AqiSensorScriptScrapped, len(sensors))
	layout := "2006-01-02T15:04:05.999999999Z"
	for _, sensor := range sensors {
		sensorTime, err := time.Parse(layout, sensor.CreatedAt)

		if err != nil {
			lib.LogError("filterDeadSensors", "failed to parse time from sensor.CreatedAt", err)
		}

		diffInHours := sensorTime.Sub(time.Now().UTC()).Hours()
		if diffInHours > float64(-1*allowedDiffInHours) {
			aliveSensors = append(aliveSensors, sensor)
		}
	}

	return aliveSensors
}
