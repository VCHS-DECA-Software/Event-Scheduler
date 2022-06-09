package links

import (
	"main/components/common"
	"main/components/globals"
	"main/components/object"

	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"
)

//Link, a type that holds an association between object.Object's
//the type parameters that distinguish it must be kept in order
type Link[F, T any] struct {
	ID   string `storm:"id"`
	From string
	To   string
}

func NewLink[F, T any](from *object.Object[F], to *object.Object[T]) (*Link[F, T], error) {
	link := &Link[F, T]{
		ID:   uuid.NewV4().String(),
		From: from.ID,
		To:   to.ID,
	}
	err := globals.DB.Save(&link)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func GetLinks[F, T any](id string) ([]Link[F, T], error) {
	var links []Link[F, T]
	query := globals.DB.Select(
		q.Or(
			q.Eq("From", id),
			q.Eq("To", id),
		),
	)
	err := query.Find(&links)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func FindLink[F, T any](from, to string) (Link[F, T], error) {
	var links []Link[F, T]
	query := globals.DB.Select(
		q.And(
			q.Eq("From", from),
			q.Eq("To", to),
		),
	)
	err := query.Find(&links)
	if err != nil {
		return Link[F, T]{}, err
	}
	return links[0], nil
}

//Find finds all other objects corresponding to the
//ID ("id") and type ("C") given
func Search[F, T, C any](id string) ([]object.Object[C], error) {
	return SearchWith[F, T, C](id, q.True())
}

//Find find one object with ID of the given "target", corresponding to the
//ID ("id") and type ("C") given
func Find[F, T, C any](id string, target string) (object.Object[C], error) {
	results, err := SearchWith[F, T, C](id, q.Eq("ID", target))
	if err != nil {
		return object.Object[C]{}, err
	}
	return results[0], nil
}

/* this method is unnecessarily difficult to use however here's a quick rundown
- type parameters "F" and "T" are the parameters you specified to a "Link"
or the parameters you called "NewLink" with (note: order matters)
- type parameter "C" is one of the parameters "F" or "T", depending on which
one you wish to search for. thus the id passed to the method must be
the id of the corresponding type.
- the "match" parameter defines the struct field matcher to run on the
filtered results (type of object.Object) */
func SearchWith[F, T, C any](id string, match q.Matcher) ([]object.Object[C], error) {
	links, err := GetLinks[F, T](id)
	if err != nil {
		return nil, err
	}

	var result []object.Object[C]
	for _, l := range links {
		other := l.From
		if l.From == id {
			other = l.To
		}

		var t object.Object[C]
		err := common.FromID(&t, other)
		if err != nil {
			return nil, err
		}

		matches, err := match.Match(other)
		if err != nil {
			return nil, err
		}

		if matches {
			result = append(result, t)
		}
	}

	return result, nil
}

func (l *Link[F, T]) Delete() error {
	return globals.DB.DeleteStruct(l)
}
