package user

import (
	"air-quality-notifyer/internal/constants"
	"air-quality-notifyer/internal/exception"
	"air-quality-notifyer/internal/helper"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	NotFound = errors.New("user not found")
)

const (
	selectUserIDByTelegramIDQuery = "SELECT id FROM users WHERE telegram_id = $1"

	deleteObservedDistrictsQuery = "DELETE FROM users_observed_districts WHERE user_id = $1"
	insertObservedDistrictQuery  = "INSERT INTO users_observed_districts (user_id, district_id) VALUES ($1, $2)"

	deleteObservedSensorsQuery    = "DELETE FROM users_observed_sensors WHERE user_id = $1"
	selectSensorIDByAPIIDQuery    = "SELECT id FROM sensors WHERE api_id = $1 ORDER BY id LIMIT 1"
	insertSensorByAPIIDQuery      = "INSERT INTO sensors (api_id) VALUES ($1) RETURNING id"
	insertObservedSensorByIDQuery = "INSERT INTO users_observed_sensors (user_id, sensor_id) VALUES ($1, $2)"
)

type Interface interface {
	FindById(ctx context.Context, id int64) (*User, error)
	Register(ctx context.Context, user User) error
	GetAllIds(ctx context.Context) ([]int64, error)
	GetAllIdsByOperatingMode(ctx context.Context, mode constants.ModeType) ([]int64, error)
	GetAllNames(ctx context.Context) ([]string, error)
	GetObservedDistrictIdsByOperatingMode(ctx context.Context, mode constants.ModeType) (map[int64][]int64, error)
	GetObservedSensorAPIIdsByOperatingMode(ctx context.Context, mode constants.ModeType) (map[int64][]int64, error)
	DeleteUserById(ctx context.Context, id int64) error
	SetOperatingMode(ctx context.Context, tgId int64, mode constants.ModeType) error
	SetObservedDistricts(ctx context.Context, tgId int64, districtIDs []int64) error
	SetObservedSensorsByAPIIds(ctx context.Context, tgId int64, sensorAPIIDs []int64) error
}

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindById(ctx context.Context, id int64) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, `
		SELECT id, username, telegram_id, operating_mode
		FROM users WHERE telegram_id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("find user by telegram id %d: %w", id, NotFound)
		}

		return nil, fmt.Errorf("find user by telegram id %d: %w", id, err)
	}

	return &user, nil
}

func (r *Repository) Register(ctx context.Context, user User) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO users (username, telegram_id)
		VALUES (:username, :telegram_id)
	`, user)

	if err != nil {
		zap.L().Error("Failed to insert user", zap.Error(err))
		return fmt.Errorf("register user: %w", err)
	}

	return nil
}

func (r *Repository) GetAllIds(ctx context.Context) ([]int64, error) {
	var ids []int64
	err := r.db.SelectContext(ctx, &ids, "SELECT telegram_id FROM users")

	if err != nil {
		return nil, fmt.Errorf("get all user ids: %w", err)
	}

	return ids, nil
}

func (r *Repository) GetAllIdsByOperatingMode(ctx context.Context, mode constants.ModeType) ([]int64, error) {
	var ids []int64
	err := r.db.SelectContext(ctx, &ids, "SELECT telegram_id FROM users WHERE operating_mode = $1", mode)
	if err != nil {
		return nil, fmt.Errorf("get user ids by operating mode %d: %w", mode, err)
	}

	return ids, nil
}

type userObservedDistrict struct {
	TelegramID int64 `db:"telegram_id"`
	DistrictID int64 `db:"district_id"`
}

func (row userObservedDistrict) telegramID() int64 {
	return row.TelegramID
}

func (row userObservedDistrict) observedID() int64 {
	return row.DistrictID
}

func (r *Repository) GetObservedDistrictIdsByOperatingMode(ctx context.Context, mode constants.ModeType) (map[int64][]int64, error) {
	var rows []userObservedDistrict
	err := r.db.SelectContext(ctx, &rows, `
		SELECT
			u.telegram_id AS telegram_id,
			uod.district_id AS district_id
		FROM users u
		JOIN users_observed_districts uod ON u.id = uod.user_id
		WHERE u.operating_mode = $1
	`, mode)
	if err != nil {
		return nil, fmt.Errorf("get observed district ids by operating mode %d: %w", mode, err)
	}

	return groupObservedIDs(rows), nil
}

type userObservedSensor struct {
	TelegramID  int64 `db:"telegram_id"`
	SensorAPIID int64 `db:"sensor_api_id"`
}

func (row userObservedSensor) telegramID() int64 {
	return row.TelegramID
}

func (row userObservedSensor) observedID() int64 {
	return row.SensorAPIID
}

func (r *Repository) GetObservedSensorAPIIdsByOperatingMode(ctx context.Context, mode constants.ModeType) (map[int64][]int64, error) {
	var rows []userObservedSensor
	err := r.db.SelectContext(ctx, &rows, `
		SELECT
			u.telegram_id AS telegram_id,
			s.api_id AS sensor_api_id
		FROM users u
		JOIN users_observed_sensors uos ON u.id = uos.user_id
		JOIN sensors s ON s.id = uos.sensor_id
		WHERE u.operating_mode = $1
	`, mode)
	if err != nil {
		return nil, fmt.Errorf("get observed sensor api ids by operating mode %d: %w", mode, err)
	}

	return groupObservedIDs(rows), nil
}

func (r *Repository) GetAllNames(ctx context.Context) ([]string, error) {
	var names []string
	err := r.db.SelectContext(ctx, &names, "SELECT username FROM users")

	if err != nil {
		return nil, fmt.Errorf("get all user names: %w", err)
	}

	return names, nil
}

func (r *Repository) DeleteUserById(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE telegram_id = $1`, id)

	if err != nil {
		return fmt.Errorf("delete user by telegram id %d: %w", id, err)
	}

	return nil
}

func (r *Repository) SetOperatingMode(ctx context.Context, tgId int64, mode constants.ModeType) error {
	if !helper.IsValidMode(mode) {
		err := exception.InvalidOperatingMode
		zap.L().Error("Setting mode", zap.Error(err))
		return err
	}

	_, err := r.db.ExecContext(ctx, "UPDATE users SET operating_mode = $1 WHERE telegram_id = $2", mode, tgId)

	if err != nil {
		zap.L().Error("Failed to set operating mode", zap.Error(err))
		return fmt.Errorf("set operating mode for user %d: %w", tgId, err)
	}

	return nil
}

func (r *Repository) SetObservedDistricts(ctx context.Context, tgId int64, districtIDs []int64) (err error) {
	transaction, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin set observed districts transaction: %w", err)
	}
	defer rollbackOnError(transaction, &err)

	userID, err := userIDByTelegramID(ctx, transaction, tgId)
	if err != nil {
		return err
	}

	if err = replaceObservedDistricts(ctx, transaction, userID, districtIDs); err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("commit set observed districts transaction: %w", err)
	}

	return nil
}

func (r *Repository) SetObservedSensorsByAPIIds(ctx context.Context, tgId int64, sensorAPIIDs []int64) (err error) {
	transaction, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin set observed sensors transaction: %w", err)
	}
	defer rollbackOnError(transaction, &err)

	userID, err := userIDByTelegramID(ctx, transaction, tgId)
	if err != nil {
		return err
	}

	if err = replaceObservedSensors(ctx, transaction, userID, sensorAPIIDs); err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("commit set observed sensors transaction: %w", err)
	}

	return nil
}

type observedRow interface {
	telegramID() int64
	observedID() int64
}

func groupObservedIDs[T observedRow](rows []T) map[int64][]int64 {
	grouped := make(map[int64][]int64)
	for _, row := range rows {
		grouped[row.telegramID()] = append(grouped[row.telegramID()], row.observedID())
	}

	return grouped
}

func rollbackOnError(transaction *sqlx.Tx, err *error) {
	if *err == nil {
		return
	}

	if rollbackErr := transaction.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
		zap.L().Error("failed to rollback transaction", zap.Error(rollbackErr))
	}
}

func userIDByTelegramID(ctx context.Context, transaction *sqlx.Tx, tgID int64) (int64, error) {
	var userID int64
	err := transaction.GetContext(ctx, &userID, selectUserIDByTelegramIDQuery, tgID)
	if err != nil {
		return 0, fmt.Errorf("get user id by telegram id %d: %w", tgID, err)
	}

	return userID, nil
}

func replaceObservedDistricts(ctx context.Context, transaction *sqlx.Tx, userID int64, districtIDs []int64) error {
	_, err := transaction.ExecContext(ctx, deleteObservedDistrictsQuery, userID)
	if err != nil {
		return fmt.Errorf("delete observed districts for user %d: %w", userID, err)
	}

	for _, districtID := range districtIDs {
		_, err = transaction.ExecContext(ctx, insertObservedDistrictQuery, userID, districtID)
		if err != nil {
			return fmt.Errorf("insert observed district %d for user %d: %w", districtID, userID, err)
		}
	}

	return nil
}

func replaceObservedSensors(ctx context.Context, transaction *sqlx.Tx, userID int64, sensorAPIIDs []int64) error {
	_, err := transaction.ExecContext(ctx, deleteObservedSensorsQuery, userID)
	if err != nil {
		return fmt.Errorf("delete observed sensors for user %d: %w", userID, err)
	}

	for _, sensorAPIID := range sensorAPIIDs {
		sensorID, err := sensorIDByAPIID(ctx, transaction, sensorAPIID)
		if err != nil {
			return err
		}

		_, err = transaction.ExecContext(ctx, insertObservedSensorByIDQuery, userID, sensorID)
		if err != nil {
			return fmt.Errorf("insert observed sensor %d for user %d: %w", sensorID, userID, err)
		}
	}

	return nil
}

func sensorIDByAPIID(ctx context.Context, transaction *sqlx.Tx, sensorAPIID int64) (int64, error) {
	var sensorID int64
	err := transaction.GetContext(ctx, &sensorID, selectSensorIDByAPIIDQuery, sensorAPIID)
	if errors.Is(err, sql.ErrNoRows) {
		err = transaction.GetContext(ctx, &sensorID, insertSensorByAPIIDQuery, sensorAPIID)
	}
	if err != nil {
		return 0, fmt.Errorf("get or create sensor by api id %d: %w", sensorAPIID, err)
	}

	return sensorID, nil
}
