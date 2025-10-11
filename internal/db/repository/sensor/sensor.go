package sensor

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type Interface interface {
	GetAllApiIds() ([]int64, error)
	GetSensorByApiId(id int64) (*Sensor, error)
	SaveSensor(sensor Sensor) error
	EvictSensor(id int64) error
	GetSensorsByDistrictId(id int64) ([]Sensor, error)
}

func (r *Repository) GetSensorByApiId(id int64) (*Sensor, error) {
	var sensor Sensor
	err := r.db.Get(&sensor, "SELECT * FROM sensors WHERE api_id = $1", id)
	if err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (r *Repository) SaveSensor(sensor Sensor) error {
	_, err := r.db.NamedExec(`
		INSERT INTO sensors (api_id, district_id, address, lat, lon, created_at)
		VALUES (:api_id, :district_id, :address, :lat, :lon, :created_at)
	`, sensor)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllApiIds() ([]int64, error) {
	var ids []int64
	err := r.db.Select(&ids, "SELECT api_id FROM sensors")

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Repository) EvictSensor(sensorApiId int64) error {
	_, err := r.db.Exec("DELETE FROM sensors WHERE api_id = $1", sensorApiId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetSensorsByDistrictId(id int64) ([]Sensor, error) {
	var sensors []Sensor
	err := r.db.Select(&sensors, `
		SELECT
		    s.id AS id,
			s.api_id AS api_id,
			s.district_id AS district_id,
			s.address AS address,
			s.lat AS lat,
			s.lon AS lon,
			d.name AS "district.name"
		FROM sensors AS s
		LEFT JOIN districts d on d.id = s.district_id
		WHERE district_id = $1
    `, id)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}
