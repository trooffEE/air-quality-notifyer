package repository

import (
	"air-quality-notifyer/internal/db/exceptions"
	"air-quality-notifyer/internal/db/models"
	"github.com/jmoiron/sqlx"
	"log"
)

type SensorRepository struct {
	db *sqlx.DB
}

func NewSensorRepository(db *sqlx.DB) *SensorRepository {
	return &SensorRepository{db: db}
}

type SensorRepositoryType interface {
	GetAllApiIds() (*[]int64, error)
	GetSensorByApiId(id int64) (*models.AirqualitySensor, error)
	SaveSensor(sensor models.AirqualitySensor) error
	EvictSensor(id int64) error
}

func (r *SensorRepository) GetSensorByApiId(id int64) (*models.AirqualitySensor, error) {
	var sensor models.AirqualitySensor
	err := r.db.Get(&sensor, "SELECT * FROM sensors WHERE api_id = $1", id)
	if err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (r *SensorRepository) SaveSensor(sensor models.AirqualitySensor) error {
	_, err := r.db.NamedExec(`
		INSERT INTO sensors (api_id, district_id, address, lat, lon)
		VALUES (:api_id, :district_id, :address, :lat, :lon)
	`, sensor)

	if err != nil {
		log.Printf("%w\n", err)
		return err
	}

	return nil
}

func (r *SensorRepository) GetAllApiIds() (*[]int64, error) {
	var ids []int64
	err := r.db.Select(&ids, "SELECT api_id FROM sensors")

	if err != nil {
		return nil, exceptions.ErrInternalDBError
	}

	return &ids, nil
}

func (r *SensorRepository) EvictSensor(sensorApiId int64) error {
	_, err := r.db.Exec("DELETE FROM sensors WHERE api_id = $1", sensorApiId)
	if err != nil {
		return exceptions.ErrInternalDBError
	}

	return nil
}
