package links

import (
	"errors"
	"main/components/db"
	"main/components/globals"
	"main/components/types"

	"github.com/asdine/storm/q"
	uuid "github.com/satori/go.uuid"
)

//Link, a type that holds an association between object.Object's
//the type parameters that distinguish it must be kept in order
type Link struct {
	ID   string `storm:"id"`
	From string
	To   string
	Type int
}

type HasAnnotations interface {
	GetID() string
	GetType() types.Type
}

func NewLink(from, to HasAnnotations) (*Link, error) {
	fromExists, err := db.Has[any](from.GetID())
	if !fromExists {
		return nil, errors.New("the \"from\" id does not exist in database")
	}
	if err != nil {
		return nil, err
	}

	toExists, err := db.Has[any](to.GetID())
	if !toExists {
		return nil, errors.New("the \"to\" id does not exist in database")
	}
	if err != nil {
		return nil, err
	}

	link := &Link{
		ID:   uuid.NewV4().String(),
		From: from.GetID(),
		To:   to.GetID(),
		Type: types.Composite(from.GetType(), to.GetType()),
	}

	err = db.Save(&link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func Select(t1 HasAnnotations) {
	links := globals.DB.Select(
		q.Eq(),
	)
}

func GetLinks(id string) ([]Link, error) {
	var links []Link
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

func FindLink(from, to string) (Link, error) {
	var links []Link
	query := globals.DB.Select(
		q.And(
			q.Eq("From", from),
			q.Eq("To", to),
		),
	)
	err := query.Find(&links)
	if err != nil {
		return Link{}, err
	}
	return links[0], nil
}

//Search finds all other objects corresponding to the
//ID ("id") and type ("C") given
func Search[T any](id string) ([]T, error) {
	return SearchWith[T](id, q.True())
}

//Find find one object with ID of the given "target", corresponding to the
//ID ("id") and type ("C") given
func Find[T any](id string, target string) (T, error) {
	results, err := SearchWith[T](id, q.Eq("ID", target))
	if err != nil {
		var r T
		return r, err
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
func SearchWith[T any](id string, match q.Matcher) ([]T, error) {
	links, err := GetLinks(id)
	if err != nil {
		return nil, err
	}

	var result []T
	for _, l := range links {
		other := l.From
		if l.From == id {
			other = l.To
		}

		found, err := db.Get[T](other)
		if err != nil {
			return nil, err
		}

		matches, err := match.Match(other)
		if err != nil {
			return nil, err
		}

		if matches {
			result = append(result, found)
		}
	}

	return result, nil
}

func (l *Link) Delete() error {
	return globals.DB.DeleteStruct(l)
}
