package sensor

import (
	"sync"
)

type Data struct {
	SensorId       int64
	Date           string
	SDS_P2         float64
	SDS_P1         float64
	Temperature    float64
	Humidity       int64
	Pressure       int64
	District       string
	AQIPM25        float64
	AQIPM10        float64
	DangerLevel    string
	DangerColor    string
	AdditionalInfo string
}

func NewSensorsData() [][]Data {
	return [][]Data{}
}

var PMLevelAirMap = []struct {
	PM25Low   float64
	PM25High  float64
	PM10Low   float64
	PM10High  float64
	IndexLow  float64
	IndexHigh float64
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

func (s *Data) calculateAQI(wg *sync.WaitGroup) {
	defer wg.Done()
	for _, pm := range PMLevelAirMap {
		if s.SDS_P1 >= pm.PM10Low && s.SDS_P1 < pm.PM10High {
			s.AQIPM10 = ((pm.IndexHigh-pm.IndexLow)/(pm.PM10High-pm.PM10Low))*(s.SDS_P1-pm.PM10Low) + pm.IndexLow
			s.DangerLevel = pm.Name
			s.DangerColor = pm.Color
		}
		if s.SDS_P2 >= pm.PM25Low && s.SDS_P2 < pm.PM25High {
			s.AQIPM25 = ((pm.IndexHigh-pm.IndexLow)/(pm.PM25High-pm.PM25Low))*(s.SDS_P2-pm.PM25Low) + pm.IndexLow
			s.DangerLevel = pm.Name
			s.DangerColor = pm.Color
		}
	}
}
