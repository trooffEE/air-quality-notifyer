package repository

import (
	"air-quality-notifyer/pkg/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	UserNotFound           = errors.New("User not found")
)

type UserRepositoryInterface interface {
	FindById(id int64) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	Create(user entity.User) (*entity.User, error)
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindById(id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, UserNotFound
	}

	if err != nil {
		fmt.Printf("%w", err)
		return nil, ErrInternalServerError
	}

	return &user, ErrInternalServerError
}

func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user)
	if err != nil {
		fmt.Printf("%w", err)
		return nil, ErrInternalServerError
	}

	return &user, ErrInternalServerError
}

func (r *UserRepository) Create(user entity.User) (*entity.User, error) {
	_, err := r.db.Exec("INSERT INTO users (id, username) VALUES (?, ?)", user.Id, user.Username)
	if err != nil {
		return nil, ErrInternalServerError
	}
	return &user, nil
}

func (*UserRepository) Update(ctx context.Context, id string) (*entity.User, error) {
	return nil, ErrInternalServerError
}

func (*UserRepository) Delete(ctx context.Context, id string) error {
	return ErrInternalServerError
}
