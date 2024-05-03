package sensor

import "sync"

type SensorData struct {
	Sensor_Id          int
	Date               string
	SDS_P2             float32
	SDS_P1             float32
	BME280_temperature float32
	BME280_humidity    float32
	BME280_pressure    interface{}
}

type SensorDataHandled struct {
	SensorData
	District string
	AQIPM25  int
	AQIPM10  int
}

func NewSensorsData() [][]SensorDataHandled {
	return [][]SensorDataHandled{}
}

var S = []struct {
	PM25Low   float64
	PM25High  float64
	PM10Low   int
	PM10High  int
	IndexLow  int
	IndexHigh int
	Color     string
	Name      string
}{
	{
		PM25Low:   0,
		PM25High:  12,
		PM10Low:   0,
		PM10High:  54,
		IndexLow:  0,
		IndexHigh: 50,
		Color:     "#50ccaa",
		Name:      "Хорошо",
	},
	{
		PM25Low:   12.1,
		PM25High:  35.4,
		PM10Low:   55,
		PM10High:  154,
		IndexLow:  51,
		IndexHigh: 100,
		Color:     "#f0e641",
		Name:      "Приемлемо",
	},
	{
		PM25Low:   35.5,
		PM25High:  55.4,
		PM10Low:   155,
		PM10High:  254,
		IndexLow:  101,
		IndexHigh: 150,
		Color:     "#fa912a",
		Name:      "Плохо",
	},
	{
		PM25Low:   55.5,
		PM25High:  150.4,
		PM10Low:   255,
		PM10High:  354,
		IndexLow:  151,
		IndexHigh: 200,
		Color:     "#ff5050",
		Name:      "Вредно",
	},
	{
		PM25Low:   150.5,
		PM25High:  250.4,
		PM10Low:   355,
		PM10High:  424,
		IndexLow:  201,
		IndexHigh: 300,
		Color:     "#8f3f97",
		Name:      "Очень вредно",
	},
	{
		PM25Low:   250.5,
		PM25High:  350.4,
		PM10Low:   425,
		PM10High:  504,
		IndexLow:  301,
		IndexHigh: 400,
		Color:     "#960032",
		Name:      "Чрезвычайно опасно",
	},
	{
		PM25Low:   350.5,
		PM25High:  500.4,
		PM10Low:   505,
		PM10High:  604,
		IndexLow:  401,
		IndexHigh: 500,
		Color:     "#960032",
		Name:      "Чрезвычайно опасно",
	},
}

func (s *SensorDataHandled) calculateAQI(wg *sync.WaitGroup) {
	defer wg.Done()
	// AQI = ((AQI_high - AQI_low) / (Conc_high - Conc_low)) * (Conc_measured - Conc_low) + AQI_low
	s.AQIPM25 = 1
	s.AQIPM25 = 1
}
