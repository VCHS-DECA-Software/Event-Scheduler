package users

import (
	"go/types"
	"main/components/db"
	"main/components/encryption"

	uuid "github.com/satori/go.uuid"
)

const (
	STUDENT = iota
	JUDGE
	ADMIN
)

type Account struct {
	ID         string `storm:"id"`
	Username   string
	Password   string
	types.Type `storm:"inline"`
}

func NewAccount(a Account) (Account, error) {
	hashedPassword, err := encryption.HashPassword(a.Password)
	if err != nil {
		return Account{}, err
	}

	account := Account{
		ID:       uuid.NewV4().String(),
		Username: a.Username,
		Password: hashedPassword,
		Type:     a.Type,
	}

	err = db.Save(&account)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}
