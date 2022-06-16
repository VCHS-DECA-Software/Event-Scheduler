package db

import (
	"main/components/globals"

	"github.com/asdine/storm"
)

/* these are some common DB methods with
stricter typing implemented on them */

func Get[T any](id string) (T, error) {
	var found T
	err := globals.DB.One("ID", id, &found)
	return found, err
}

func Has[T any](id string) (bool, error) {
	_, err := Get[T](id)
	if err == storm.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func Save[T any](object *T) error {
	return globals.DB.Save(object)
}

func Update[T any](object *T) error {
	return globals.DB.Update(object)
}

func Delete[T any](object *T) error {
	return globals.DB.DeleteStruct(object)
}
