package sensor

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
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
	metalploshadka = "metalploshadka"
	lesnayaPolyana = "lesnayaPolyana"
)

// TODO вынести в отдельный пакет
type requestDone struct {
	_sensors []Data
}

func (r requestDone) notifyChangesInSensors() {
	NotifyChangesInSensors(r._sensors)
}

type districtsWithFallback struct {
	name        string
	fallbackIds []int64
}

var districtNames = map[string]string{
	center:         "Центральный",
	kirovskii:      "Кировский",
	circus:         `"Цирк"`,
	boulevard:      "Бульвар",
	yuzhinii:       "Южный",
	metalploshadka: "Металлплощадка",
	lesnayaPolyana: "Лесная Поляна",
}

var districts map[int64]districtsWithFallback = map[int64]districtsWithFallback{
	7: districtsWithFallback{
		boulevard,
		[]int64{},
	},
	11: districtsWithFallback{
		lesnayaPolyana,
		[]int64{},
	},
	20: districtsWithFallback{
		metalploshadka,
		[]int64{53},
	},
	40: districtsWithFallback{
		center,
		[]int64{39, 48},
	},
	47: districtsWithFallback{
		kirovskii,
		[]int64{},
	},
	59: districtsWithFallback{
		yuzhinii,
		[]int64{51, 56},
	},
	71: districtsWithFallback{
		circus,
		[]int64{},
	},
}

func PingForSensorsDataOnceIn(cronString string) {
	sensors := NewSensorsData()
	c := cron.New()
	_, err := c.AddFunc(cronString, func() { fetchSensors(sensors).notifyChangesInSensors() })
	if err != nil {
		log.Panic(err)
	}
	c.Start()
	select {}
}

func fetchSensors(sensors []Data) requestDone {
	respChan := make(chan Data, len(districts))
	var wg sync.WaitGroup
	wg.Add(len(districts))

	for key, value := range districts {
		go fetchSensorById(&wg, respChan, key, value)
	}

	wg.Wait()
	close(respChan)

	for resp := range respChan {
		sensors = append(sensors, resp)
	}

	return requestDone{sensors}
}

func fetchSensorById(wg *sync.WaitGroup, resChan chan Data, id int64, districtInfo districtsWithFallback) {
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
		log.Printf("\nfetchSensorById http status code %d for \"%s\" District with %d, revoking with fallback sensor", res.StatusCode, districts[id].name, id)
		for fallbackId := range districtInfo.fallbackIds {
			// Done to make sure that "fallback go routines" won't fill data for "main go routines" districts
			go fetchSensorById(
				wg,
				resChan,
				int64(fallbackId),
				districtsWithFallback{districtInfo.name, []int64{}},
			)
		}
		return
	}

	// TODO Переписать так, чтобы этот метод использовался только если нужны все данные, а не только данные о последнем сенсоре
	// TODO нас интересует последний ответ в массиве, он является актуальным для текущего часа
	resChan <- richSensorData(fetchedSensorData[len(fetchedSensorData)-1], districtInfo, id)
	wg.Done()
	defer res.Body.Close()
}

func richSensorData(
	fetchedSensorData Data,
	districtInfoRelatedToFetchedSensorData districtsWithFallback,
	id int64,
) Data {
	sensorData := fetchedSensorData

	sensorData.District = districtInfoRelatedToFetchedSensorData.name
	sensorData.Id = id
	if sensorData.Humidity >= 90 {
		sensorData.AdditionalInfo = "Высокая влажность. Показания PM могут быть не корректны\n"
	}
	if sensorData.Temperature < -60 {
		sensorData.AdditionalInfo += fmt.Sprintf("Датчики температуры в районе %s не исправен!\n", sensorData.District)
	}
	sensorData.getInformationAboutAQI()

	return sensorData
}
