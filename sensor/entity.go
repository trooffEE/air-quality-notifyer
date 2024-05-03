package sensor

type SensorData struct {
	Sensor_Id          int
	Date               string
	SDS_P2             float32
	SDS_P1             float32
	BME280_temperature float32
	BME280_humidity    float32
	BME280_pressure    interface{}
}

type SensorDataWithDistrict struct {
	SensorData
	district string
}

func NewSensorsData() [][]SensorData {
	return [][]SensorData{}
}
