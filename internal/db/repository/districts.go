package repository

import (
	"air-quality-notifyer/internal/db/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DistrictRepository struct {
	db *sqlx.DB
}

func NewDistrictRepository(db *sqlx.DB) *DistrictRepository {
	return &DistrictRepository{db: db}
}

type DistrictRepositoryInterface interface {
	GetAllDistricts() ([]models.District, error)
	GetAssociatedDistrictIdByCoords(x, y float64) *models.DistrictSensor
}

func (r *DistrictRepository) GetAssociatedDistrictIdByCoords(x, y float64) *models.DistrictSensor {
	var sensorDistrict models.DistrictSensor
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

func (r *DistrictRepository) GetAllDistricts() ([]models.District, error) {
	var districts []models.District
	err := r.db.Select(&districts, "SELECT d.id, d.name FROM districts AS d")
	return districts, err
}
