package users

import (
	"fmt"
	"main/components/db"
	"main/components/encryption"
	"main/components/links"
)

func IssueToken[T Type](id string) (string, error) {
	t, key, err := encryption.IssueToken()
	if err != nil {
		return "", err
	}

	_, err = links.NewLink[T, encryption.Key](id, key)
	if err != nil {
		return "", err
	}

	return t, nil
}

func Authenticate[T Type](id string, password string) (string, error) {
	a, err := db.Get[Account[T]](id)
	if err != nil {
		return "", err
	}

	if encryption.CheckPasswordHash(password, a.Password) {
		return "", fmt.Errorf("invalid password")
	}

	return IssueToken[T](id)
}

func Validate[T Type](id string, token string) (bool, error) {
	return encryption.CheckToken(token)
}
