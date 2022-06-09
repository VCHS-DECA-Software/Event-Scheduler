package methods

import (
	"main/components/groupings"
	"main/components/object"
	"main/components/users"
)

/* when making associations, make sure
to use the order [account, grouping]
for consistency */

func CreateTeam(
	account *object.Object[users.Account[users.Student]],
	name string,
) error {
	return groupings.CreateGrouping(
		account, groupings.Hierarchy[groupings.Team]{
			Name: name,
		},
	)
}

func CreateEvent(
	account *object.Object[users.Account[users.Admin]],
	name string, values groupings.Event,
) error {
	return groupings.CreateGrouping(
		account, groupings.Hierarchy[groupings.Event]{
			Name: name,
			Data: groupings.Event{
				Description: values.Description,
				Location:    values.Location,
				EventType:   values.EventType,
				StartTime:   values.StartTime,
				EndTime:     values.EndTime,
			},
		},
	)
}
