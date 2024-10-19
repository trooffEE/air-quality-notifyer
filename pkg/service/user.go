package service

import (
	"air-quality-notifyer/pkg/entity"
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

func (ur *UserService) IsNewUser(id int64) bool {
	_, err := ur.repo.FindById(id)
	if errors.Is(repo.UserNotFound, err) {
		return true
	}

	if err != nil {
		fmt.Println(err)
	}

	return false
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
