package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/dioxine/grpc-pb/internal/models"
	interfaces "github.com/dioxine/grpc-pb/pkg/v1"
)

type UseCase struct {
	repo interfaces.RepoInterface
}

func New(repo interfaces.RepoInterface) interfaces.UseCaseInterface {
	return &UseCase{repo}
}

func (uc *UseCase) Create(user models.User) (models.User, error) {

	if _, err := uc.repo.GetByEmail(user.Email); err == nil {
		return models.User{}, errors.New("the email is already associated with another user")
	}

	// create new user
	user, err := uc.repo.Create(user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (uc *UseCase) Get(id string, email string, password string) (models.User, error) {

	user, err := uc.repo.Get(id, email, password)
	if err != nil {
		log.Printf("no such user with the id: %s or email: %s supplied", id, email)
		return models.User{}, err
	}
	return user, nil
}

func (uc *UseCase) Update(updateUser models.User) error {
	user, err := uc.repo.Get(updateUser.Id, "", "")

	if err != nil {
		return err
	}

	if user.Email != updateUser.Email {
		return fmt.Errorf("email cannot be changed")
	}

	if err := uc.repo.Update(updateUser); err != nil {
		log.Println("failed to update user")
		return err
	} else {
		log.Println("updated user:", updateUser)
	}

	return nil
}

func (uc *UseCase) Delete(id string) error {

	if user, err := uc.repo.Get(id, "", ""); err != nil {
		return err
	} else {
		log.Println("user to be deleted:", user)
	}

	if err := uc.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
