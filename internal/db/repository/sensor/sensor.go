package sensor

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type Interface interface {
	GetAllApiIds(ctx context.Context) ([]int64, error)
	GetSensorByApiId(ctx context.Context, id int64) (*Sensor, error)
	SaveSensor(ctx context.Context, sensor Sensor) error
	EvictSensor(ctx context.Context, id int64) error
	GetSensorsByDistrictId(ctx context.Context, id int64) ([]Sensor, error)
}

func (r *Repository) GetSensorByApiId(ctx context.Context, id int64) (*Sensor, error) {
	var sensor Sensor
	err := r.db.GetContext(ctx, &sensor, "SELECT * FROM sensors WHERE api_id = $1", id)
	if err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (r *Repository) SaveSensor(ctx context.Context, sensor Sensor) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO sensors (api_id, district_id, address, lat, lon, created_at)
		VALUES (:api_id, :district_id, :address, :lat, :lon, :created_at)
	`, sensor)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllApiIds(ctx context.Context) ([]int64, error) {
	var ids []int64
	err := r.db.SelectContext(ctx, &ids, "SELECT api_id FROM sensors")

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Repository) EvictSensor(ctx context.Context, sensorApiId int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sensors WHERE api_id = $1", sensorApiId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetSensorsByDistrictId(ctx context.Context, id int64) ([]Sensor, error) {
	var sensors []Sensor
	err := r.db.SelectContext(ctx, &sensors, `
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
