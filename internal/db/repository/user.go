package repository

import (
	"air-quality-notifyer/internal/db/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	UserNotFound = errors.New("User not found")
)

type UserRepositoryInterface interface {
	FindById(id int64) (*models.User, error)
	Register(user models.User) error
	GetAllIds() ([]int64, error)
	GetAllNames() ([]string, error)
	DeleteUserById(id int64) error
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindById(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, `
		SELECT id, username, telegram_id, operating_mode
		FROM users WHERE telegram_id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Register(user models.User) error {
	_, err := r.db.NamedExec(`
		INSERT INTO users (username, telegram_id)
		VALUES (:username, :telegram_id)
	`, user)

	if err != nil {
		zap.L().Error("Failed to insert user", zap.Error(err))
		return err
	}

	return nil
}

func (r *UserRepository) GetAllIds() ([]int64, error) {
	var ids []int64
	err := r.db.Select(&ids, "SELECT telegram_id FROM users")

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *UserRepository) GetAllNames() ([]string, error) {
	var names []string
	err := r.db.Select(&names, "SELECT username FROM users")

	if err != nil {
		return nil, err
	}

	return names, nil
}

func (r *UserRepository) DeleteUserById(id int64) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE telegram_id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}
