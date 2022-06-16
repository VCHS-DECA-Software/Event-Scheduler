package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/asdine/storm"
)

func Check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func InitializeDB(t *testing.T, name string) *storm.DB {
	db, err := storm.Open(fmt.Sprintf("%v.db", name))
	if err != nil {
		t.Error(err)
	}
	return db
}

func CleanupDB(t *testing.T, name string, db *storm.DB) {
	Check(t, db.Close())
	Check(t, os.Remove(fmt.Sprintf("%v.db", name)))
}
