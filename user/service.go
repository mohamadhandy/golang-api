package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
}

// hanya di package user saja boleh panggil service service ini.
type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	fmt.Println("repo: ", repository, "service: ", &service{repository})
	printService := &service{repository}
	fmt.Println("isi repo: ", &repository, "isi service: ", *printService)
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"
	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

// mapping struct input ke struct User
// simpan struct User melalui repository
