package repo

import (
	"log"

	"github.com/dioxine/grpc-pb/internal/models"
	interfaces "github.com/dioxine/grpc-pb/pkg/v1"
	"github.com/pocketbase/pocketbase/daos"
	pbmodels "github.com/pocketbase/pocketbase/models"
)

type Repo struct {
	db *daos.Dao
}

func New(db *daos.Dao) interfaces.RepoInterface {
	return &Repo{db}
}
func (repo *Repo) Create(user models.User) (models.User, error) {
	collection, err := repo.db.FindCollectionByNameOrId("users")
	if err != nil {
		log.Println(err)
	}
	record := pbmodels.NewRecord(collection)
	// set individual fields
	// or bulk load with record.Load(map[string]any{...}, except Password)
	record.Set("emailVisibility", false)
	record.Set("username", &user.Username)
	record.Set("name", &user.Name)
	record.Set("email", &user.Email)
	record.SetPassword(user.Password)
	record.SetVerified(true)

	if err := repo.db.SaveRecord(record); err != nil {
		return models.User{}, err
	}

	return models.User{
		Id:       record.GetId(),
		Username: record.Username(),
		Name:     record.GetString("name"),
		Email:    record.Email(),
		TokenKey: record.TokenKey(),
	}, nil
}

func (repo *Repo) Get(id string, email string, password string) (models.User, error) {
	if email != "" {
		record, err := repo.db.FindAuthRecordByEmail("users", email)
		if err != nil {
			return models.User{}, err
		}
		return models.User{
			Id:           record.GetId(),
			Username:     record.Username(),
			Name:         record.GetString("name"),
			Email:        record.Email(),
			PasswordIsOk: record.ValidatePassword(password),
			TokenKey:     record.TokenKey(),
		}, nil
	}

	record, err := repo.db.FindRecordById("users", id)
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		Id:           record.GetId(),
		Username:     record.Username(),
		Name:         record.GetString("name"),
		Email:        record.Email(),
		PasswordIsOk: record.ValidatePassword(password),
		TokenKey:     record.TokenKey(),
	}, nil
}

func (repo *Repo) Update(user models.User) error {
	record, err := repo.db.FindRecordById("users", user.Id)
	if err != nil {
		log.Println(err)
	}
	// set individual fields
	// or bulk load with record.Load(map[string]any{...})
	record.Set("emailVisibility", false)
	record.Set("username", &user.Username)
	record.Set("name", &user.Name)
	record.Set("email", &user.Email)
	record.SetPassword(user.Password)

	if err := repo.db.SaveRecord(record); err != nil {
		log.Println(err)
	}
	return nil
}
func (repo *Repo) Delete(id string) error {
	record, err := repo.db.FindRecordById("users", id)
	if err != nil {
		log.Println(err)
	}
	if err := repo.db.DeleteRecord(record); err != nil {
		log.Println(err)
	}
	return nil
}

// GetByEmail
//
// This function returns the user instance which is
// saved on the DB and returns to the usecase

func (repo *Repo) GetByEmail(email string) (models.User, error) {
	record, err := repo.db.FindFirstRecordByData("users", "email", email)
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		Email: record.GetString("email"),
		Name:  record.GetString("name"),
	}, nil
}
