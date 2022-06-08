package methods

import (
	"main/components/dbmanager"
	"main/components/links"
	"main/components/users"
	uuid "main/vendor/github.com/satori/go.uuid"
)

func CreateTeam(student *users.Account[users.Student], name string) error {
	return dbmanager.Save(&links.Team{
		ID:       uuid.NewV4().String(),
		Name:     name,
		OwnerID:  student.ID,
		Students: links.NewLink(-1),
		Events:   links.NewLink(2),
	})
}
