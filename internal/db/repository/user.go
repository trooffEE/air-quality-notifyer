package repository

import (
	"air-quality-notifyer/internal/db/exceptions"
	"air-quality-notifyer/internal/db/models"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
)

type UserRepositoryInterface interface {
	FindById(id int64) (*models.User, error)
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindById(id int64) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, exceptions.UserNotFound
	}

	if err != nil {
		log.Printf("%w\n", err)
		return nil, exceptions.ErrInternalDBError
	}

	return &user, exceptions.ErrInternalDBError
}
