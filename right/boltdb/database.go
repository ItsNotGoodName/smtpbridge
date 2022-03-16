package boltdb

import (
	"context"
	"log"

	"github.com/asdine/storm/v3"
)

type Database struct {
	db *storm.DB
}

func NewDatabase(path string) Database {
	db, err := storm.Open(path)
	if err != nil {
		log.Fatalln("bolt.NewDatabase: could not open database:", err)
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
	log.Println("boltdb.Database.Run: closed database")
	done <- struct{}{}
}
