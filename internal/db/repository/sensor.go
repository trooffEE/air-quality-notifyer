package repository

import "github.com/jmoiron/sqlx"

type SensorRepository struct {
	db *sqlx.DB
}

func NewSensorRepository(db *sqlx.DB) *SensorRepository {
	return &SensorRepository{db: db}
}

type SensorRepositoryType interface {
	GetDistrict(id int64) error
}

func (r *SensorRepository) GetDistrict(_ int64) error {
	return nil
}
