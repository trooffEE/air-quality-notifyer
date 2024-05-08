package sensor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	center         = "center"
	kirovskii      = "kirovskii"
	circus         = "circus"
	boulevard      = "boulevard"
	yuzhinii       = "yuzhinii"
	metalposhadka  = "metalposhadka"
	lesnayaPolyana = "lesnayaPolyana"
)

type DistrictsWithFallback struct {
	name        string
	fallbackIds []int64
}

var DistrictNames = map[string]string{
	center:         "Центральный",
	kirovskii:      "Кировский",
	circus:         `"Цирк"`,
	boulevard:      "Бульвар",
	yuzhinii:       "Южный",
	metalposhadka:  "Металлплощадка",
	lesnayaPolyana: "Лесная Поляна",
}

var Districts map[int64]DistrictsWithFallback = map[int64]DistrictsWithFallback{
	7: DistrictsWithFallback{
		boulevard,
		[]int64{},
	},
	11: DistrictsWithFallback{
		lesnayaPolyana,
		[]int64{},
	},
	20: DistrictsWithFallback{
		metalposhadka,
		[]int64{53},
	},
	40: DistrictsWithFallback{
		center,
		[]int64{39, 48},
	},
	47: DistrictsWithFallback{
		kirovskii,
		[]int64{},
	},
	59: DistrictsWithFallback{
		yuzhinii,
		[]int64{51, 56},
	},
	71: DistrictsWithFallback{
		circus,
		[]int64{},
	},
}

func FetchSensorsData(sensors *[][]Data) {
	respChan := make(chan []Data, len(Districts))
	var wg sync.WaitGroup
	wg.Add(len(Districts))

	for key, value := range Districts {
		go fetchSensorById(&wg, respChan, key, value)
	}

	wg.Wait()
	close(respChan)

	for resp := range respChan {
		*sensors = append(*sensors, resp)
	}

	ChangesInAPIAppearedChannel <- *sensors
}

func fetchSensorById(wg *sync.WaitGroup, resChan chan []Data, id int64, districtInfo DistrictsWithFallback) {
	var fetchedSensorData []Data

	res, err := http.Post(
		fmt.Sprintf("https://airkemerovo.ru/api/sensor/archive/%d/1", id),
		"application/json",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(res.Body).Decode(&fetchedSensorData)
	if err != nil {
		log.Println("Something went wrong on decoding JSON from API step")
	}

	if res.StatusCode != http.StatusOK || fetchedSensorData == nil || len(fetchedSensorData) == 0 {
		log.Printf("\nfetchSensorById http status code %d for \"%s\" district with %d, revoking with fallback sensor", res.StatusCode, Districts[id].name, id)
		for fallbackId := range districtInfo.fallbackIds {
			// Done to make sure that "fallback go routines" won't fill data for "main go routines" districts
			go fetchSensorById(
				wg,
				resChan,
				int64(fallbackId),
				DistrictsWithFallback{districtInfo.name, []int64{}},
			)
		}
		return
	}

	richSensorData(fetchedSensorData, districtInfo, id)

	resChan <- fetchedSensorData
	wg.Done()
	defer res.Body.Close()
}

func richSensorData(
	fetchedSensorData []Data,
	districtInfoRelatedToFetchedSensorData DistrictsWithFallback,
	id int64,
) []Data {
	var wg sync.WaitGroup
	wg.Add(len(fetchedSensorData))

	for i := range fetchedSensorData {
		sensorData := &fetchedSensorData[i]
		sensorData.SensorId = id
		sensorData.District = districtInfoRelatedToFetchedSensorData.name
		if sensorData.Humidity >= 90 {
			sensorData.AdditionalInfo = "Высокая влажность. Показания PM могут быть не корректны\n"
		}
		if sensorData.Temperature < -60 {
			sensorData.AdditionalInfo += fmt.Sprintf("Датчики температуры в районе %s не исправен!\n", sensorData.District)
		}
		go sensorData.calculateAQI(&wg)
	}
	wg.Wait()

	return fetchedSensorData
}
