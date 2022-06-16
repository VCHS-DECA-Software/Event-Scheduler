package groupings

import (
	"fmt"
	"main/components/db"
	"main/components/links"
	"main/components/users"

	"github.com/asdine/storm"
)

func CreateGrouping[A users.Type, G GroupingType](
	account users.Account[A],
	grouping Hierarchy[G],
) error {
	grouping.Owner = account.ID
	err := db.Save(&grouping)
	if err != nil {
		return err
	}
	_, err = links.NewLink[A, G](account.ID, grouping.ID)
	return err
}

func JoinGrouping[A users.Type, G GroupingType](
	account users.Account[A],
	id string,
) error {
	_, err := links.Find[A, G, G](account.ID, id)
	if err == nil {
		return fmt.Errorf("account has already joined grouping")
	}
	if err.Error() != storm.ErrNotFound.Error() {
		return err
	}

	event, err := db.Get[Hierarchy[G]](id)
	if err != nil {
		return err
	}

	_, err = links.NewLink[A, G](account.ID, event.ID)
	return err
}

func LeaveGrouping[A users.Type, G GroupingType](
	account users.Account[A],
	id string,
) error {
	link, err := links.FindLink[A, G](account.ID, id)
	if err != nil {
		return err
	}
	return link.Delete()
}
