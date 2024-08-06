package models

type User struct {
	Id           string `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	Email        string `db:"email" json:"email"`
	Name         string `db:"name" json:"name"`
	Password     string `db:"password" json:"password"`
	TokenKey     string `db:"token_key" json:"token_key"`
	PasswordIsOk bool
}
