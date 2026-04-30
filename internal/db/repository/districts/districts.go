package districts

import (
	"air-quality-notifyer/internal/db/repository/sensor"
	"context"
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
	GetAllDistrictsNames(ctx context.Context) ([]string, error)
	GetAllDistricts(ctx context.Context) ([]District, error)
	GetAssociatedDistrictIdByCoords(ctx context.Context, x, y float64) *sensor.DistrictSensor
}

func (r *Repository) GetAssociatedDistrictIdByCoords(ctx context.Context, x, y float64) *sensor.DistrictSensor {
	var sensorDistrict sensor.DistrictSensor
	var pointGeo = fmt.Sprintf("SRID=4326;POINT(%f %f)", x, y)
	err := r.db.GetContext(ctx, &sensorDistrict, `
		SELECT id, name
		FROM districts as d
		WHERE st_contains(d.area, $1)
	`, pointGeo)
	if err != nil {
		return nil
	}

	return &sensorDistrict
}

func (r *Repository) GetAllDistricts(ctx context.Context) ([]District, error) {
	var districts []District
	err := r.db.SelectContext(ctx, &districts, "SELECT d.id, d.name FROM districts AS d ORDER BY d.name DESC")
	return districts, err
}

func (r *Repository) GetAllDistrictsNames(ctx context.Context) ([]string, error) {
	var districts []string
	err := r.db.SelectContext(ctx, &districts, "SELECT d.name FROM districts AS d ORDER BY d.name DESC")
	return districts, err
}
