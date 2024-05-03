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

func fetchSensorById(wg *sync.WaitGroup, resChan chan []SensorData, id int, fallback DistrictsWithFallback) {
	defer wg.Done()
	var fetchedSensorData []SensorData

	res, err := http.Post(
		fmt.Sprintf("https://airkemerovo.ru/api/sensor/archive/%d/1", id),
		"application/json",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	json.NewDecoder(res.Body).Decode(&fetchedSensorData)

	if res.StatusCode != http.StatusOK || fetchedSensorData == nil || len(fetchedSensorData) == 0 {
		log.Printf("fetchSensorById http status code %d for \"%s\" district with %d, revoking with fallback sensor", res.StatusCode, SensorIdsMapForDistricts[id].name, id)
		for fallbackId := range fallback.fallbackIds {
			// Done to make sure that "fallback go routines" won't fill data for "main go routines" districts
			// TODO Think about better solution
			time.Sleep(100)
			go fetchSensorById(wg, resChan, fallbackId, DistrictsWithFallback{fallback.name, []int{}})
		}
		return
	}

	resChan <- fetchedSensorData
	defer res.Body.Close()
}

func FetchSensorsData(sensors *[][]SensorData) {
	respChan := make(chan []SensorData, len(SensorIdsMapForDistricts))
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
