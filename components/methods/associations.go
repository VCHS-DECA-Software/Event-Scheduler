package methods

import (
	"fmt"
	"main/components/groupings"
	"main/components/links"
	"main/components/users"

	"github.com/asdine/storm"
)

func JoinEvent[T users.Judge | users.Student](
	account users.Account[T], eventID string,
) error {
	_, err := links.Find[T, groupings.Event, groupings.Event](
		account.ID, eventID,
	)
	if err == nil {
		return fmt.Errorf(
			"this account has already been enrolled in the event",
		)
	}
	if err.Error() != storm.ErrNotFound.Error() {
		return err
	}

	_, err = links.NewLink[T, groupings.Event](account.ID, eventID)
	return err
}

func LeaveEvent[T users.Judge | users.Student](
	account users.Account[T], eventID string,
) error {
	link, err := links.FindLink[T, groupings.Event](account.ID, eventID)
	if err != nil {
		return err
	}
	return link.Delete()
}
