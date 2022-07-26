package boltdb

import (
	"context"
	"log"

	"github.com/asdine/storm/v3"
)

type Database struct {
	db *storm.DB
}

func NewDatabase(file string) Database {
	db, err := storm.Open(file)
	if err != nil {
		log.Fatalln("boltdb.NewDatabase: could not open database:", err)
	}

	return Database{
		db: db,
	}
}

func (d Database) Close() error {
	return d.db.Close()
}

func (d Database) Run(ctx context.Context, done chan<- struct{}) {
	<-ctx.Done()
	d.Close()
	done <- struct{}{}
}

func count(tx storm.Node, inf interface{}) (int, error) {
	count, err := tx.Count(inf)
	if err != nil {
		if err == storm.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}
