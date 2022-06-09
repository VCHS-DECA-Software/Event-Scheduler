package users

import (
	"main/components/globals"
	"main/components/object"
)

func Initialize() error {
	err := globals.DB.Init(&object.Object[Account[Student]]{})
	if err != nil {
		return err
	}
	err = globals.DB.Init(&object.Object[Account[Judge]]{})
	if err != nil {
		return err
	}
	return globals.DB.Init(&object.Object[Account[Admin]]{})
}
