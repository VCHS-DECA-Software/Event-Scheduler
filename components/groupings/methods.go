package groupings

import (
	"errors"
	"main/components/globals"
	"main/components/links"
	"main/components/object"
	"main/components/users"

	"github.com/asdine/storm"
)

func CreateGrouping[A users.Type, G any](
	account *object.Object[users.Account[A]],
	grouping Hierarchy[G],
) error {
	grouping.Owner = account.ID
	group, err := object.NewObject(grouping)
	if err != nil {
		return err
	}
	_, err = links.NewLink(account, group)
	return err
}

func JoinGrouping[A users.Type, G any](
	account *object.Object[users.Account[A]],
	id string,
) error {
	_, err := links.Find[
		users.Account[A], Hierarchy[G], Hierarchy[G],
	](account.ID, id)
	if err == nil {
		return errors.New("account has already joined groupipng")
	}
	if err.Error() != storm.ErrNotFound.Error() {
		return err
	}

	var event object.Object[Hierarchy[G]]
	err = globals.DB.One("ID", id, &event)
	if err != nil {
		return err
	}

	_, err = links.NewLink(account, &event)
	return err
}

func LeaveGrouping[A users.Type, G any](
	account *object.Object[users.Account[A]],
	id string,
) error {
	link, err := links.FindLink[
		users.Account[A], Hierarchy[G],
	](account.ID, id)
	if err != nil {
		return err
	}
	return link.Delete()
}
