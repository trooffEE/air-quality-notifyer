package districts

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type Interface interface {
	GetAllDistricts() ([]District, error)
	GetAssociatedDistrictIdByCoords(x, y float64) *sensor.DistrictSensor
}

func (r *Repository) GetAssociatedDistrictIdByCoords(x, y float64) *sensor.DistrictSensor {
	var sensorDistrict sensor.DistrictSensor
	var pointGeo = fmt.Sprintf("SRID=4326;POINT(%f %f)", x, y)
	err := r.db.Get(&sensorDistrict, `
		SELECT id, name
		FROM districts as d
		WHERE st_contains(d.area, $1)
	`, pointGeo)
	if err != nil {
		return nil
	}

	return &sensorDistrict
}

func (r *Repository) GetAllDistricts() ([]District, error) {
	var districts []District
	err := r.db.Select(&districts, "SELECT d.id, d.name FROM districts AS d ORDER BY d.name DESC")
	return districts, err
}
