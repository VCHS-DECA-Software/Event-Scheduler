package users

import (
	"main/components/db"
	"main/components/encryption"

	uuid "github.com/satori/go.uuid"
)

const (
	STUDENT = iota
	JUDGE
	ADMIN
)

/* these empty structs exist for differentiation, please take a look
at docs/database.md a more detailed explanation */
type Student struct {
}
type Judge struct {
}
type Admin struct {
}

type Type interface {
	Student | Judge | Admin
}

type Account[T Type] struct {
	ID       string `storm:"id"`
	Username string
	Password string
}

func NewAccount[T Type](a Account[T]) (Account[T], error) {
	hashedPassword, err := encryption.HashPassword(a.Password)
	if err != nil {
		return Account[T]{}, err
	}

	account := Account[T]{
		ID:       uuid.NewV4().String(),
		Username: a.Username,
		Password: hashedPassword,
	}

	err = db.Save(&account)
	if err != nil {
		return Account[T]{}, err
	}

	return account, nil
}
