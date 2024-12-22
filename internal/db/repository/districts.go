package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DistrictRepository struct {
	db *sqlx.DB
}

func NewDistrictRepository(db *sqlx.DB) *DistrictRepository {
	return &DistrictRepository{db: db}
}

type DistrictRepositoryType interface {
	GetAssociatedDistrictIdByCoords(x, y float64) int64
}

func (r *DistrictRepository) GetAssociatedDistrictIdByCoords(x, y float64) int64 {
	var id int64
	var pointGeo = fmt.Sprintf("SRID=4326;POINT(%f %f)", x, y)
	err := r.db.Get(&id, "SELECT id as area FROM districts where st_contains(area, $1)", pointGeo)
	if err != nil {
		return -1
	}

	return id
}
