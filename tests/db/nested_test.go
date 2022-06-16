package db_test

import (
	"main/tests/utils"
	"testing"
)

type ParentComposite struct {
	ID    string `storm:"id"`
	Child `storm:"inline"`
}

type ParentNested struct {
	ID    string `storm:"id"`
	Child Child  `storm:"inline"`
}

type Child struct {
	Name string `storm:"index"`
}

func TestCompositeQuery(t *testing.T) {
	db := utils.InitializeDB(t, "composite_query")
	defer utils.CleanupDB(t, "composite_query", db)

	utils.Check(t, db.Init(&ParentComposite{}))
	utils.Check(t, db.Save(&ParentComposite{
		ID: "id1",
		Child: Child{
			Name: "name1",
		},
	}))

	var parent ParentComposite
	utils.Check(t, db.One("Name", "name1", &parent))
	if parent.Name != "name1" {
		t.Errorf("composite query via inline struct failed")
	}
}

func TestNestedQuery(t *testing.T) {
	db := utils.InitializeDB(t, "nested_query")
	defer utils.CleanupDB(t, "nested_query", db)

	utils.Check(t, db.Init(&ParentNested{}))
	utils.Check(t, db.Save(&ParentNested{
		ID: "id1",
		Child: Child{
			Name: "name1",
		},
	}))

	var parent ParentNested
	utils.Check(t, db.One("Name", "name1", &parent))
	if parent.Child.Name == "name1" {
		t.Errorf("nested query via inline struct should not have succeeded")
	}
}
