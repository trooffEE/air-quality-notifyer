package user

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	NotFound = errors.New("user not found")
)

type Interface interface {
	FindById(id int64) (*User, error)
	Register(user User) error
	GetAllIds() ([]int64, error)
	GetAllNames() ([]string, error)
	DeleteUserById(id int64) error
	//SetMode(mode int) error
}

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindById(id int64) (*User, error) {
	var user User
	err := r.db.Get(&user, `
		SELECT id, username, telegram_id, operating_mode
		FROM users WHERE telegram_id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) Register(user User) error {
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

func (r *Repository) GetAllIds() ([]int64, error) {
	var ids []int64
	err := r.db.Select(&ids, "SELECT telegram_id FROM users")

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Repository) GetAllNames() ([]string, error) {
	var names []string
	err := r.db.Select(&names, "SELECT username FROM users")

	if err != nil {
		return nil, err
	}

	return names, nil
}

func (r *Repository) DeleteUserById(id int64) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE telegram_id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SetMode() error {

	return nil
}
