package usecase

import (
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
		log.Println("the email is already associated with another user")
		return models.User{}, err
	}

	// create new user
	user, err := uc.repo.Create(user)
	if err != nil {
		log.Println("failed to create new user")
		return models.User{}, err
	} else {
		log.Println("new user created:", user)
	}
	return user, nil
}

func (uc *UseCase) Get(id string) (models.User, error) {

	user, err := uc.repo.Get(id)
	if err != nil {
		log.Println("no such user with the id supplied")
		return models.User{}, err
	} else {
		log.Println("user found:", user)
	}
	return user, nil
}

func (uc *UseCase) Update(updateUser models.User) error {
	user, err := uc.repo.Get(updateUser.Id)

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

	if user, err := uc.repo.Get(id); err != nil {
		return err
	} else {
		log.Println("user to be deleted:", user)
	}

	if err := uc.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
