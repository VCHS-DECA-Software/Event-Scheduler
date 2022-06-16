package users

import (
	"main/components/globals"
)

func Initialize() error {
	err := globals.DB.Init(&Account[Student]{})
	if err != nil {
		return err
	}
	err = globals.DB.Init(&Account[Judge]{})
	if err != nil {
		return err
	}
	return globals.DB.Init(&Account[Admin]{})
}
