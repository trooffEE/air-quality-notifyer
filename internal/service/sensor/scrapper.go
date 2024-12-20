package sensor

import (
	"bufio"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

type AirqualitySensorScriptScrapped struct {
	Id      int64   `json:"id"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

var setSensorsStringStart, setSensorsStringEnd = "setSensors('", "');"

func scrapSensorData() []AirqualitySensorScriptScrapped {
	res, err := http.Get("https://airkemerovo.ru")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
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

	var sensors []AirqualitySensorScriptScrapped
	for scanner.Scan() {
		scriptLine := scanner.Text()

		if strings.Index(scriptLine, setSensorsStringStart) != -1 {
			startJsonIndex := strings.Index(scriptLine, setSensorsStringStart) + len(setSensorsStringStart)
			endJsonIndex := strings.Index(scriptLine, setSensorsStringEnd)

			jsonString := scriptLine[startJsonIndex:endJsonIndex]
			err := json.Unmarshal([]byte(jsonString), &sensors)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return sensors
}
