package models

type AirqualitySensor struct {
	Id         int64   `db:"id"`
	ApiId      int64   `db:"api_id"`
	DistrictId int64   `db:"district_id"`
	Address    string  `db:"address"`
	Lat        float64 `db:"lat"`
	Lon        float64 `db:"lon"`
	District   struct {
		Id   int64  `db:"id"`
		Name string `db:"name"`
	}
}
