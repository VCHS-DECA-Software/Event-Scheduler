package methods

import (
	"errors"
	"main/components/globals"
	"main/components/groupings"
	"main/components/links"
	"main/components/object"
	"main/components/users"

	"github.com/asdine/storm"
)

func JoinEvent[
	T users.Judge | users.Student,
](account *object.Object[users.Account[T]], eventID string) error {
	_, err := links.Find[
		users.Account[T],
		groupings.Event,
		groupings.Event,
	](account.ID, eventID)
	if err == nil {
		return errors.New("this account has already been enrolled in the event")
	}
	if err.Error() != storm.ErrNotFound.Error() {
		return err
	}

	var event object.Object[groupings.Event]
	err = globals.DB.One("ID", eventID, &event)
	if err != nil {
		return err
	}

	_, err = links.NewLink(account, &event)
	return err
}

func LeaveEvent[
	T users.Judge | users.Student,
](account *object.Object[users.Account[T]], eventID string) error {
	link, err := links.FindLink[
		users.Account[T],
		groupings.Event,
	](account.ID, eventID)
	if err != nil {
		return err
	}
	return link.Delete()
}
