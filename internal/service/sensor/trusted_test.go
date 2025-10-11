package sensor

import (
	"air-quality-notifyer/internal/service/sensor/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteSyncAirqualitySensorList struct {
	suite.Suite
}

func TestSuites(t *testing.T) {
	suite.Run(t, new(SuiteSyncAirqualitySensorList))
}

func (s *SuiteSyncAirqualitySensorList) TestAddSensor() {
	t := s.T()
	t.Parallel()
	sensorsList := model.SyncSensorsList{}

	sensorsList.AddSensor(model.Sensor{Id: 1})
	sensorsList.AddSensor(model.Sensor{Id: 2})
	sensorsList.AddSensor(model.Sensor{Id: 3})

	assert.Equal(t, sensorsList.List, []model.Sensor{{Id: 1}, {Id: 2}, {Id: 3}})
}

func (s *SuiteSyncAirqualitySensorList) TestSortAqi() {
	t := s.T()
	t.Parallel()
	sensorsList := model.SyncSensorsList{}

	sensorsList.AddSensor(model.Sensor{Aqi: 99})
	sensorsList.AddSensor(model.Sensor{Aqi: 1})
	sensorsList.AddSensor(model.Sensor{Aqi: 33})
	sensorsList.AddSensor(model.Sensor{Aqi: 33})
	sensorsList.AddSensor(model.Sensor{Aqi: 66})

	expectedResult := []model.Sensor{{Aqi: 1}, {Aqi: 33}, {Aqi: 33}, {Aqi: 66}, {Aqi: 99}}

	assert.NotEqual(t, sensorsList.List, expectedResult)
	sensorsList.SortByAqi()
	assert.Equal(t, sensorsList.List, expectedResult)
}

func (s *SuiteSyncAirqualitySensorList) TestGetTrustedAqiSensor_empty() {
	t := s.T()
	t.Parallel()
	sensorsList := model.SyncSensorsList{}

	trustedSensor := sensorsList.GetSensor()
	assert.Nil(t, trustedSensor)
}

func (s *SuiteSyncAirqualitySensorList) TestGetTrustedAqiSensor_odd() {
	t := s.T()
	t.Parallel()
	sensorsList := model.SyncSensorsList{}

	sensorsList.AddSensor(model.Sensor{Aqi: 10})
	sensorsList.AddSensor(model.Sensor{Aqi: 100})
	sensorsList.AddSensor(model.Sensor{Aqi: 50})

	trustedSensor := sensorsList.GetSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, model.Sensor{Aqi: 50})
}

func (s *SuiteSyncAirqualitySensorList) TestGetTrustedAqiSensor_even() {
	t := s.T()
	t.Parallel()
	sensorsList := model.SyncSensorsList{}

	sensorsList.AddSensor(model.Sensor{Aqi: 10})
	sensorsList.AddSensor(model.Sensor{Aqi: 100})
	sensorsList.AddSensor(model.Sensor{Aqi: 50})
	sensorsList.AddSensor(model.Sensor{Aqi: 33})

	trustedSensor := sensorsList.GetSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, model.Sensor{Aqi: 50})
}

func (s *SuiteSyncAirqualitySensorList) TestGetTrustedAqiSensor_one() {
	t := s.T()
	t.Parallel()
	sensorsList := model.SyncSensorsList{}

	sensorsList.AddSensor(model.Sensor{Aqi: 10})

	trustedSensor := sensorsList.GetSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, model.Sensor{Aqi: 10})
}
