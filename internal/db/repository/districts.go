package repository

import (
	"air-quality-notifyer/internal/db/exceptions"
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
	GetAssociatedDistrictIdByCoords(x, y float64) (int64, error)
}

func (r *DistrictRepository) GetAssociatedDistrictIdByCoords(x, y float64) (int64, error) {
	var id int64
	var pointGeo = fmt.Sprintf("SRID=4326;POINT(%f %f)", x, y)
	err := r.db.Get(&id, "SELECT id as area FROM districts where st_contains(area, $1)", pointGeo)
	if err != nil {
		return -1, exceptions.ErrInternalDBError
	}

	return id, nil
}
