package sensor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
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
	fallbackIds []int
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

var SensorIdsMapForDistricts map[int]DistrictsWithFallback = map[int]DistrictsWithFallback{
	7: DistrictsWithFallback{
		boulevard,
		[]int{},
	},
	11: DistrictsWithFallback{
		lesnayaPolyana,
		[]int{},
	},
	20: DistrictsWithFallback{
		metalposhadka,
		[]int{53},
	},
	40: DistrictsWithFallback{
		center,
		[]int{39, 48},
	},
	47: DistrictsWithFallback{
		kirovskii,
		[]int{},
	},
	56: DistrictsWithFallback{
		yuzhinii,
		[]int{51, 59},
	},
	71: DistrictsWithFallback{
		circus,
		[]int{},
	},
}

func FetchSensorsData(sensors *[][]SensorDataHandled) {
	respChan := make(chan []SensorDataHandled, len(SensorIdsMapForDistricts))
	var wg sync.WaitGroup
	wg.Add(len(SensorIdsMapForDistricts))

	for key, value := range SensorIdsMapForDistricts {
		go fetchSensorById(&wg, respChan, key, value)
	}

	wg.Wait()
	close(respChan)

	for resp := range respChan {
		*sensors = append(*sensors, resp)
	}
}

func fetchSensorById(wg *sync.WaitGroup, resChan chan []SensorDataHandled, id int, districtInfo DistrictsWithFallback) {
	defer wg.Done()
	var fetchedSensorData []SensorDataHandled

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
		log.Printf("\nfetchSensorById http status code %d for \"%s\" district with %d, revoking with fallback sensor", res.StatusCode, SensorIdsMapForDistricts[id].name, id)
		for fallbackId := range districtInfo.fallbackIds {
			// Done to make sure that "fallback go routines" won't fill data for "main go routines" districts
			// TODO Think about better solution
			time.Sleep(100)
			go fetchSensorById(wg, resChan, fallbackId, DistrictsWithFallback{districtInfo.name, []int{}})
		}
		return
	}

	handleFetchSensorData(fetchedSensorData, districtInfo)

	resChan <- fetchedSensorData
	defer res.Body.Close()
}

func handleFetchSensorData(fetchedSensorData []SensorDataHandled, districtInfoRelatedToFetchedSensorData DistrictsWithFallback) {
	var wg sync.WaitGroup
	wg.Add(len(fetchedSensorData))

	for i := range fetchedSensorData {
		sensorData := &fetchedSensorData[i]
		sensorData.District = districtInfoRelatedToFetchedSensorData.name
		go sensorData.calculateAQI(&wg)
	}

	wg.Wait()
}
