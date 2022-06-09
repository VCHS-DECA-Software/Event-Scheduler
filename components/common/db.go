package common

import "main/components/globals"

func FromID[T any](value *T, id string) error {
	return globals.DB.One("ID", id, value)
}
