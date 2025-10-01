package sensor

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SuiteSyncAirqualitySensorList struct {
	suite.Suite
}

func (s *SuiteSyncAirqualitySensorList) TestAddSensor() {
	t := s.T()
	t.Parallel()
	sensorsList := SyncSensors{}

	sensorsList.addSensor(Sensor{Id: 1})
	sensorsList.addSensor(Sensor{Id: 2})
	sensorsList.addSensor(Sensor{Id: 3})

	assert.Equal(t, sensorsList.list, []Sensor{{Id: 1}, {Id: 2}, {Id: 3}})
}

func (s *SuiteSyncAirqualitySensorList) TestSortAqi() {
	t := s.T()
	t.Parallel()
	sensorsList := SyncSensors{}

	sensorsList.addSensor(Sensor{Aqi: 99})
	sensorsList.addSensor(Sensor{Aqi: 1})
	sensorsList.addSensor(Sensor{Aqi: 33})
	sensorsList.addSensor(Sensor{Aqi: 33})
	sensorsList.addSensor(Sensor{Aqi: 66})

	expectedResult := []Sensor{{Aqi: 1}, {Aqi: 33}, {Aqi: 33}, {Aqi: 66}, {Aqi: 99}}

	assert.NotEqual(t, sensorsList.list, expectedResult)
	sensorsList.sortByAqi()
	assert.Equal(t, sensorsList.list, expectedResult)
}

func (s *SuiteSyncAirqualitySensorList) TestGetTrustedAqiSensor_empty() {
	t := s.T()
	t.Parallel()
	sensorsList := SyncSensors{}

	trustedSensor := sensorsList.getTrustedSensor()
	assert.Nil(t, trustedSensor)
}

func (s *SuiteSyncAirqualitySensorList) TestGetTrustedAqiSensor_odd() {
	t := s.T()
	t.Parallel()
	sensorsList := SyncSensors{}

	sensorsList.addSensor(Sensor{Aqi: 10})
	sensorsList.addSensor(Sensor{Aqi: 100})
	sensorsList.addSensor(Sensor{Aqi: 50})

	trustedSensor := sensorsList.getTrustedSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, Sensor{Aqi: 50})
}

func (s *SuiteSyncAirqualitySensorList) TestGetTrustedAqiSensor_even() {
	t := s.T()
	t.Parallel()
	sensorsList := SyncSensors{}

	sensorsList.addSensor(Sensor{Aqi: 10})
	sensorsList.addSensor(Sensor{Aqi: 100})
	sensorsList.addSensor(Sensor{Aqi: 50})
	sensorsList.addSensor(Sensor{Aqi: 33})

	trustedSensor := sensorsList.getTrustedSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, Sensor{Aqi: 50})
}

func (s *SuiteSyncAirqualitySensorList) TestGetTrustedAqiSensor_one() {
	t := s.T()
	t.Parallel()
	sensorsList := SyncSensors{}

	sensorsList.addSensor(Sensor{Aqi: 10})

	trustedSensor := sensorsList.getTrustedSensor()
	assert.NotNil(t, trustedSensor)
	assert.Equal(t, *trustedSensor, Sensor{Aqi: 10})
}
