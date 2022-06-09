package db_test

import (
	"log"
	"testing"

	"github.com/asdine/storm"
)

type Generic[T any] struct {
	ID   string `storm:"id"`
	Data T
}

type GenericUnused[T any] struct {
	ID string `storm:"id"`
}

type Type1 struct {
	Text string
}

type Type2 struct {
	Text string
}

func TestGenericStructNested(t *testing.T) {
	db, err := storm.Open("database_nested.db")
	if err != nil {
		t.Error(err)
	}

	err = db.Init(&Generic[any]{})
	if err != nil {
		t.Error(err)
	}

	err = db.Save(&Generic[Type1]{
		ID: "id1",
		Data: Type1{
			Text: "text1",
		},
	})
	if err != nil {
		t.Error(err)
	}
	err = db.Save(&Generic[Type2]{
		ID: "id2",
		Data: Type2{
			Text: "text2",
		},
	})
	if err != nil {
		t.Error(err)
	}

	var type1 Generic[Type1]
	err = db.One("ID", "id1", &type1)
	if err != nil {
		t.Error(err)
	}
	if type1.ID != "id1" {
		t.Error("type1 ID mismatch")
	}
	if type1.Data.Text != "text1" {
		t.Error("type1 nested values mismatch")
	}

	var type2 Generic[Type2]
	err = db.One("ID", "id2", &type2)
	if err != nil {
		t.Error(err)
	}
	if type2.ID != "id2" {
		t.Error("type2 ID mismatch")
	}
	if type2.Data.Text != "text2" {
		t.Error("type2 nested values mismatch")
	}

	var shouldFail Generic[Type1]
	err = db.One("ID", "id2", &shouldFail)
	if err == nil {
		t.Error("non existent object should not be resolved")
	} else if err.Error() != storm.ErrNotFound.Error() {
		t.Error(err)
	}
	log.Println(shouldFail)
}

func TestGenericUnusedStructNested(t *testing.T) {
	db, err := storm.Open("database_nested_unused.db")
	if err != nil {
		t.Error(err)
	}

	err = db.Init(&GenericUnused[any]{})
	if err != nil {
		t.Error(err)
	}

	err = db.Save(&GenericUnused[Type1]{
		ID: "id1",
	})
	if err != nil {
		t.Error(err)
	}
	err = db.Save(&GenericUnused[Type2]{
		ID: "id2",
	})
	if err != nil {
		t.Error(err)
	}

	var type1 GenericUnused[Type1]
	err = db.One("ID", "id1", &type1)
	if err != nil {
		t.Error(err)
	}
	if type1.ID != "id1" {
		t.Error("type1 ID mismatch")
	}

	var type2 GenericUnused[Type2]
	err = db.One("ID", "id2", &type2)
	if err != nil {
		t.Error(err)
	}
	if type2.ID != "id2" {
		t.Error("type2 ID mismatch")
	}

	var shouldFail GenericUnused[Type1]
	err = db.One("ID", "id2", &shouldFail)
	if err == nil {
		t.Error("non existent object should not be resolved")
	} else if err.Error() != storm.ErrNotFound.Error() {
		t.Error(err)
	}
	log.Println(shouldFail)
}
