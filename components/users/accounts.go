package users

import (
	"main/components/dbmanager"
	"main/components/encryption"

	uuid "github.com/satori/go.uuid"
)

/*
	make an overarching "Account" struct that implements logic for
	each of the different users, since Judges, Admins and Students
	are not linked together, there is no need for them to be in
	different tables. One can simply add a "type" field to the account
	to distinguish between different user types. That way DRY can be
	utilized.

	furthermore, account types exist in some sort of hierarchy, admins
	have the most amount of access while students have the least amount
	of access. One should be able to implement a system that can quantify
	that access and use it inside each api method that performs linking
	and logic.
*/

const (
	STUDENT = iota
	JUDGE
	ADMIN
)

type Account[T any] struct {
	ID       string `storm:"id"`
	Type     int
	Username string
	Password string
	Data     T
}

func (a *Account[T]) Save(data T) error {
	id := uuid.NewV4().String()
	hashedPassword, err := encryption.HashPassword(a.Password)
	if err != nil {
		return err
	}
	err = dbmanager.Save(&Account[T]{
		ID: id, Type: a.Type,
		Username: a.Username,
		Password: hashedPassword,
		Data:     data,
	})
	return err
}

func (a *Account[T]) FromID(id string) error {
	return dbmanager.Query("ID", id, &a)
}

func (a *Account[T]) Update(updated Account[T]) error {
	hashed, err := encryption.HashPassword(updated.Password)
	if err != nil {
		return err
	}
	updated.Password = hashed
	return dbmanager.Update(&updated)
}

func (a *Account[T]) Delete() error {
	return dbmanager.Delete(a)
}
