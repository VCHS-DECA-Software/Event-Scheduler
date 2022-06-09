package object

import (
	"main/components/globals"

	uuid "github.com/satori/go.uuid"
)

type Object[T any] struct {
	ID   string `storm:"id"`
	Data T      `storm:"inline"`
}

func NewObject[T any](data T) (*Object[T], error) {
	grouping := &Object[T]{
		ID:   uuid.NewV4().String(),
		Data: data,
	}
	err := globals.DB.Save(grouping)
	if err != nil {
		return nil, err
	}
	return grouping, nil
}

func FromID[T any](id string) (*Object[T], error) {
	grouping := &Object[T]{}
	return grouping, globals.DB.One("ID", id, grouping)
}

func (o *Object[T]) Update(data T) error {
	return globals.DB.Update(&Object[T]{
		ID:   o.ID,
		Data: data,
	})
}

func (g *Object[T]) Delete() error {
	return globals.DB.DeleteStruct(g)
}
