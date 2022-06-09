package users

import (
	"main/components/encryption"
	"main/components/object"
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
	Username string
	Password string
}

func NewAccount[T Type](a Account[T]) (*object.Object[Account[T]], error) {
	hashedPassword, err := encryption.HashPassword(a.Password)
	if err != nil {
		return nil, err
	}
	return object.NewObject(
		Account[T]{
			Username: a.Username,
			Password: hashedPassword,
		},
	)
}

func (a *Account[T]) Authenticate(password string) bool {
	return encryption.CheckPasswordHash(password, a.Password)
}
