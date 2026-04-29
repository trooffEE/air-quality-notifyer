package sensor

import (
	rSensor "air-quality-notifyer/internal/db/repository/sensor"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsAliveSensorCacheKey(t *testing.T) {
	assert.True(t, isAliveSensorCacheKey("sensor:123"))
	assert.False(t, isAliveSensorCacheKey("sensor:district:1"))
	assert.False(t, isAliveSensorCacheKey("district:1"))
}

func TestAliveSensorsFromCachePayloads(t *testing.T) {
	payload, err := json.Marshal(rSensor.Sensor{
		ApiId:   123,
		Address: "Lenina 1",
		Lat:     55.3,
		Lon:     86.1,
		District: rSensor.DistrictSensor{
			Name: "Central",
		},
	})
	require.NoError(t, err)

	sensors, err := aliveSensorsFromCachePayloads([]string{string(payload)})
	require.NoError(t, err)

	require.Len(t, sensors, 1)
	assert.Equal(t, int64(123), sensors[0].APIID)
	assert.Equal(t, "Lenina 1", sensors[0].Address)
	assert.Equal(t, "Central", sensors[0].District)
}
