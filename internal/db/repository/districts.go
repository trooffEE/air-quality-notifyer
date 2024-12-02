package repository

import (
	"github.com/jmoiron/sqlx"
)

type DistrictRepository struct {
	db *sqlx.DB
}

func NewDistrictRepository(db *sqlx.DB) *DistrictRepository {
	return &DistrictRepository{db: db}
}

type DistrictRepositoryType interface {
	GetDistrict(id int64) error
}

func (r *DistrictRepository) GetDistrict(_ int64) error {
	return nil
}
