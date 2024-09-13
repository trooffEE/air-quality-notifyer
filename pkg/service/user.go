package service

import (
	"air-quality-notifyer/pkg/entity"
	errApp "air-quality-notifyer/pkg/errors"
	repo "air-quality-notifyer/pkg/repository"
	"errors"
	"fmt"
)

type UserService struct {
	repo repo.UserRepositoryInterface
}

func NewUserService(ur repo.UserRepositoryInterface) *UserService {
	return &UserService{
		repo: ur,
	}
}

func (ur *UserService) Register(id, username string) *entity.User {
	fmt.Println("Hello new User!, ")

	usr, err := ur.repo.Create(entity.User{
		Id:       id,
		Username: username,
	})

	if err != nil {
		fmt.Println("Error Appeared on creating new User record in DB")
	}

	return usr
}

func (ur *UserService) IsRegistered(id string) *entity.User {
	usr, err := ur.repo.FindById(id)
	if errors.Is(err, errApp.ErrUserNotFound) {
		fmt.Println("Error Appeared on creating new User record in DB")
		return nil
	}

	return usr
}
