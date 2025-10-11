package models

type Sensor struct {
	Id         int64          `db:"id" redis:"id"`
	ApiId      int64          `db:"api_id" redis:"api_id"`
	DistrictId int64          `db:"district_id" redis:"district_id"`
	Address    string         `db:"address" redis:"address"`
	Lat        float64        `db:"lat" redis:"lat"`
	Lon        float64        `db:"lon" redis:"lon"`
	CreatedAt  string         `db:"created_at" redis:"created_at"`
	District   DistrictSensor `redis:"district" db:"district"`
}

type DistrictSensor struct {
	Id   int64  `db:"id" redis:"id"`
	Name string `db:"name" redis:"name"`
}
