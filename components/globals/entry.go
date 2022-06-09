package globals

import "github.com/asdine/storm"

type Context struct {
	DBName string
}

func Initialize(context Context) error {
	var err error
	DB, err = storm.Open(context.DBName)
	return err
}

func Destroy() error {
	return DB.Close()
}
