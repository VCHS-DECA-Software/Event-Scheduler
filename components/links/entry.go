package links

import (
	"main/components/globals"
)

func Initialize() error {
	return globals.DB.Init(&Link{})
}
