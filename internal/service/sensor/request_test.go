package sensor

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSyncAirqualitySensorList_addSensor(t *testing.T) {
	t.Parallel()
	sensorsList := SyncAirqualitySensorList{}

	sensorsList.addSensor(AqiSensor{Id: 1})
	sensorsList.addSensor(AqiSensor{Id: 2})
	sensorsList.addSensor(AqiSensor{Id: 3})

	assert.Equal(t, sensorsList.list, []AqiSensor{{Id: 1}, {Id: 2}, {Id: 3}})
}

func TestSyncAirqualitySensorList_sortAqi(t *testing.T) {
	t.Parallel()
	sensorsList := SyncAirqualitySensorList{}

	sensorsList.addSensor(AqiSensor{Aqi: 99})
	sensorsList.addSensor(AqiSensor{Aqi: 1})
	sensorsList.addSensor(AqiSensor{Aqi: 33})
	sensorsList.addSensor(AqiSensor{Aqi: 66})

	expectedResult := []AqiSensor{{Aqi: 1}, {Aqi: 33}, {Aqi: 66}, {Aqi: 99}}

	assert.NotEqual(t, sensorsList.list, expectedResult)
	sensorsList.sortAqi()
	assert.Equal(t, sensorsList.list, expectedResult)
}

func TestSyncAirqualitySensorList_getTrustedAqiSensor_empty(t *testing.T) {
	t.Parallel()
	sensorsList := SyncAirqualitySensorList{}

	trustedSensor := sensorsList.getTrustedAqiSensor()
	assert.Nil(t, trustedSensor)
}

func TestSyncAirqualitySensorList_getTrustedAqiSensor_odd(t *testing.T) {
	t.Parallel()
	sensorsList := SyncAirqualitySensorList{}

	sensorsList.addSensor(AqiSensor{Aqi: 10})
	sensorsList.addSensor(AqiSensor{Aqi: 100})
	sensorsList.addSensor(AqiSensor{Aqi: 50})

	trustedSensor := sensorsList.getTrustedAqiSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, AqiSensor{Aqi: 50})
}

func TestSyncAirqualitySensorList_getTrustedAqiSensor_even(t *testing.T) {
	t.Parallel()
	sensorsList := SyncAirqualitySensorList{}

	sensorsList.addSensor(AqiSensor{Aqi: 10})
	sensorsList.addSensor(AqiSensor{Aqi: 100})
	sensorsList.addSensor(AqiSensor{Aqi: 50})
	sensorsList.addSensor(AqiSensor{Aqi: 33})

	trustedSensor := sensorsList.getTrustedAqiSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, AqiSensor{Aqi: 50})
}

func TestSyncAirqualitySensorList_getTrustedAqiSensor_one(t *testing.T) {
	t.Parallel()
	sensorsList := SyncAirqualitySensorList{}

	sensorsList.addSensor(AqiSensor{Aqi: 10})

	trustedSensor := sensorsList.getTrustedAqiSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, AqiSensor{Aqi: 10})
}

func TestFetchSensorById(t *testing.T) {
	t.Parallel()
	sensorId := 71

	//mockBody := io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(`{"id":%d,"city_id":1,"description":"","lat":55.34396,"lon":86.107647,"address":"\u041b\u0435\u043d\u0438\u043d\u0430, 67\u0430","floor":3,"radius":100,"source":null,"type":"stationary","last":{"aqi":21,"aqi25":21,"aqi10":12,"level":"good","color":"green","pm10":13.4,"pm25":4.95,"humidity":75.07,"temperature":-6.24,"pressure":99707.94,"pressure_hpa":997,"pressure_mmhg":748,"updated_at":"2025-01-19 14:48:51"},"archive":[{"date":"2025-01-19 14","aqi":32,"aqi25":32,"aqi10":20,"level":"good","color":"green","pm10":21.4,"pm25":7.7,"humidity":73,"temperature":-6.2,"pressure":99708,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 13","aqi":25,"aqi25":25,"aqi10":16,"level":"good","color":"green","pm10":17.3,"pm25":6.1,"humidity":71,"temperature":-6,"pressure":99709,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 12","aqi":18,"aqi25":18,"aqi10":10,"level":"good","color":"green","pm10":10.8,"pm25":4.2,"humidity":74,"temperature":-5.5,"pressure":99686,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 11","aqi":27,"aqi25":27,"aqi10":16,"level":"good","color":"green","pm10":17.7,"pm25":6.4,"humidity":77,"temperature":-5,"pressure":99643,"pressure_hpa":996,"pressure_mmhg":747},{"date":"2025-01-19 10","aqi":34,"aqi25":34,"aqi10":20,"level":"good","color":"green","pm10":21.3,"pm25":8.1,"humidity":76,"temperature":-4.9,"pressure":99630,"pressure_hpa":996,"pressure_mmhg":747},{"date":"2025-01-19 09","aqi":35,"aqi25":35,"aqi10":19,"level":"good","color":"green","pm10":20.7,"pm25":8.4,"humidity":74,"temperature":-4.3,"pressure":99665,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 08","aqi":34,"aqi25":34,"aqi10":22,"level":"good","color":"green","pm10":23.7,"pm25":8.1,"humidity":69,"temperature":-3.3,"pressure":99711,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 07","aqi":18,"aqi25":18,"aqi10":10,"level":"good","color":"green","pm10":11.3,"pm25":4.3,"humidity":61,"temperature":-1.2,"pressure":99766,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-19 06","aqi":16,"aqi25":16,"aqi10":9,"level":"good","color":"green","pm10":9.6,"pm25":3.8,"humidity":71,"temperature":-4.3,"pressure":99784,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-19 05","aqi":16,"aqi25":16,"aqi10":9,"level":"good","color":"green","pm10":9.3,"pm25":3.8,"humidity":74,"temperature":-5,"pressure":99803,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 04","aqi":20,"aqi25":20,"aqi10":11,"level":"good","color":"green","pm10":12.1,"pm25":4.9,"humidity":75,"temperature":-5.3,"pressure":99816,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 03","aqi":18,"aqi25":18,"aqi10":11,"level":"good","color":"green","pm10":11.5,"pm25":4.4,"humidity":75,"temperature":-5.7,"pressure":99825,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 02","aqi":22,"aqi25":22,"aqi10":12,"level":"good","color":"green","pm10":12.6,"pm25":5.2,"humidity":74,"temperature":-5.8,"pressure":99843,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 01","aqi":15,"aqi25":15,"aqi10":9,"level":"good","color":"green","pm10":9.4,"pm25":3.6,"humidity":75,"temperature":-5.9,"pressure":99834,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 00","aqi":15,"aqi25":15,"aqi10":9,"level":"good","color":"green","pm10":9.8,"pm25":3.5,"humidity":75,"temperature":-5.8,"pressure":99825,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-18 23","aqi":17,"aqi25":17,"aqi10":10,"level":"good","color":"green","pm10":10.7,"pm25":4.1,"humidity":76,"temperature":-5.8,"pressure":99822,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-18 22","aqi":11,"aqi25":11,"aqi10":7,"level":"good","color":"green","pm10":7.6,"pm25":2.7,"humidity":74,"temperature":-5.8,"pressure":99814,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-18 21","aqi":8,"aqi25":8,"aqi10":4,"level":"good","color":"green","pm10":4.3,"pm25":1.9,"humidity":74,"temperature":-5.7,"pressure":99786,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-18 20","aqi":8,"aqi25":8,"aqi10":5,"level":"good","color":"green","pm10":5.3,"pm25":2,"humidity":74,"temperature":-5.5,"pressure":99768,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-18 19","aqi":11,"aqi25":11,"aqi10":7,"level":"good","color":"green","pm10":7.3,"pm25":2.7,"humidity":76,"temperature":-5.1,"pressure":99755,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-18 18","aqi":21,"aqi25":21,"aqi10":14,"level":"good","color":"green","pm10":15.5,"pm25":5.1,"humidity":75,"temperature":-5,"pressure":99729,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-18 17","aqi":18,"aqi25":18,"aqi10":10,"level":"good","color":"green","pm10":11.1,"pm25":4.2,"humidity":75,"temperature":-4.9,"pressure":99691,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-18 16","aqi":25,"aqi25":25,"aqi10":16,"level":"good","color":"green","pm10":17.2,"pm25":5.9,"humidity":77,"temperature":-4.8,"pressure":99643,"pressure_hpa":996,"pressure_mmhg":747},{"date":"2025-01-18 15","aqi":30,"aqi25":30,"aqi10":17,"level":"good","color":"green","pm10":18.7,"pm25":7.1,"humidity":76,"temperature":-4.9,"pressure":99617,"pressure_hpa":996,"pressure_mmhg":747},{"date":"2025-01-18 14","aqi":56,"aqi25":56,"aqi10":39,"level":"moderate","color":"yellow","pm10":42,"pm25":14.5,"humidity":76,"temperature":-5.2,"pressure":99604,"pressure_hpa":996,"pressure_mmhg":747}]}`, sensorId))))
	expectedResponse := fmt.Sprintf(`{"id":%d,"city_id":1,"description":"","lat":55.34396,"lon":86.107647,"address":"\u041b\u0435\u043d\u0438\u043d\u0430, 67\u0430","floor":3,"radius":100,"source":null,"type":"stationary","last":{"aqi":21,"aqi25":21,"aqi10":12,"level":"good","color":"green","pm10":13.4,"pm25":4.95,"humidity":75.07,"temperature":-6.24,"pressure":99707.94,"pressure_hpa":997,"pressure_mmhg":748,"updated_at":"2025-01-19 14:48:51"},"archive":[{"date":"2025-01-19 14","aqi":32,"aqi25":32,"aqi10":20,"level":"good","color":"green","pm10":21.4,"pm25":7.7,"humidity":73,"temperature":-6.2,"pressure":99708,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 13","aqi":25,"aqi25":25,"aqi10":16,"level":"good","color":"green","pm10":17.3,"pm25":6.1,"humidity":71,"temperature":-6,"pressure":99709,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 12","aqi":18,"aqi25":18,"aqi10":10,"level":"good","color":"green","pm10":10.8,"pm25":4.2,"humidity":74,"temperature":-5.5,"pressure":99686,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 11","aqi":27,"aqi25":27,"aqi10":16,"level":"good","color":"green","pm10":17.7,"pm25":6.4,"humidity":77,"temperature":-5,"pressure":99643,"pressure_hpa":996,"pressure_mmhg":747},{"date":"2025-01-19 10","aqi":34,"aqi25":34,"aqi10":20,"level":"good","color":"green","pm10":21.3,"pm25":8.1,"humidity":76,"temperature":-4.9,"pressure":99630,"pressure_hpa":996,"pressure_mmhg":747},{"date":"2025-01-19 09","aqi":35,"aqi25":35,"aqi10":19,"level":"good","color":"green","pm10":20.7,"pm25":8.4,"humidity":74,"temperature":-4.3,"pressure":99665,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 08","aqi":34,"aqi25":34,"aqi10":22,"level":"good","color":"green","pm10":23.7,"pm25":8.1,"humidity":69,"temperature":-3.3,"pressure":99711,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-19 07","aqi":18,"aqi25":18,"aqi10":10,"level":"good","color":"green","pm10":11.3,"pm25":4.3,"humidity":61,"temperature":-1.2,"pressure":99766,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-19 06","aqi":16,"aqi25":16,"aqi10":9,"level":"good","color":"green","pm10":9.6,"pm25":3.8,"humidity":71,"temperature":-4.3,"pressure":99784,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-19 05","aqi":16,"aqi25":16,"aqi10":9,"level":"good","color":"green","pm10":9.3,"pm25":3.8,"humidity":74,"temperature":-5,"pressure":99803,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 04","aqi":20,"aqi25":20,"aqi10":11,"level":"good","color":"green","pm10":12.1,"pm25":4.9,"humidity":75,"temperature":-5.3,"pressure":99816,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 03","aqi":18,"aqi25":18,"aqi10":11,"level":"good","color":"green","pm10":11.5,"pm25":4.4,"humidity":75,"temperature":-5.7,"pressure":99825,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 02","aqi":22,"aqi25":22,"aqi10":12,"level":"good","color":"green","pm10":12.6,"pm25":5.2,"humidity":74,"temperature":-5.8,"pressure":99843,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 01","aqi":15,"aqi25":15,"aqi10":9,"level":"good","color":"green","pm10":9.4,"pm25":3.6,"humidity":75,"temperature":-5.9,"pressure":99834,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-19 00","aqi":15,"aqi25":15,"aqi10":9,"level":"good","color":"green","pm10":9.8,"pm25":3.5,"humidity":75,"temperature":-5.8,"pressure":99825,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-18 23","aqi":17,"aqi25":17,"aqi10":10,"level":"good","color":"green","pm10":10.7,"pm25":4.1,"humidity":76,"temperature":-5.8,"pressure":99822,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-18 22","aqi":11,"aqi25":11,"aqi10":7,"level":"good","color":"green","pm10":7.6,"pm25":2.7,"humidity":74,"temperature":-5.8,"pressure":99814,"pressure_hpa":998,"pressure_mmhg":749},{"date":"2025-01-18 21","aqi":8,"aqi25":8,"aqi10":4,"level":"good","color":"green","pm10":4.3,"pm25":1.9,"humidity":74,"temperature":-5.7,"pressure":99786,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-18 20","aqi":8,"aqi25":8,"aqi10":5,"level":"good","color":"green","pm10":5.3,"pm25":2,"humidity":74,"temperature":-5.5,"pressure":99768,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-18 19","aqi":11,"aqi25":11,"aqi10":7,"level":"good","color":"green","pm10":7.3,"pm25":2.7,"humidity":76,"temperature":-5.1,"pressure":99755,"pressure_hpa":998,"pressure_mmhg":748},{"date":"2025-01-18 18","aqi":21,"aqi25":21,"aqi10":14,"level":"good","color":"green","pm10":15.5,"pm25":5.1,"humidity":75,"temperature":-5,"pressure":99729,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-18 17","aqi":18,"aqi25":18,"aqi10":10,"level":"good","color":"green","pm10":11.1,"pm25":4.2,"humidity":75,"temperature":-4.9,"pressure":99691,"pressure_hpa":997,"pressure_mmhg":748},{"date":"2025-01-18 16","aqi":25,"aqi25":25,"aqi10":16,"level":"good","color":"green","pm10":17.2,"pm25":5.9,"humidity":77,"temperature":-4.8,"pressure":99643,"pressure_hpa":996,"pressure_mmhg":747},{"date":"2025-01-18 15","aqi":30,"aqi25":30,"aqi10":17,"level":"good","color":"green","pm10":18.7,"pm25":7.1,"humidity":76,"temperature":-4.9,"pressure":99617,"pressure_hpa":996,"pressure_mmhg":747},{"date":"2025-01-18 14","aqi":56,"aqi25":56,"aqi10":39,"level":"moderate","color":"yellow","pm10":42,"pm25":14.5,"humidity":76,"temperature":-5.2,"pressure":99604,"pressure_hpa":996,"pressure_mmhg":747}]}`, sensorId)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/%d", sensorId) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, expectedResponse)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	endpoint = ts.URL + "/%d"

	sensorResponse, err := fetchSensorById(int64(sensorId))
	assert.NoError(t, err)

	var lazyExpectedParsedResult AqiSensorResponse
	json.Unmarshal([]byte(expectedResponse), &lazyExpectedParsedResult)
	assert.Equal(t, sensorResponse, lazyExpectedParsedResult)
}
