package repository

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/asdine/storm/v3"
)

type Database struct {
	db *storm.DB
}

func NewDatabase(cfg *config.Config) Database {
	db, err := storm.Open(cfg.DB.DB)
	if err != nil {
		log.Fatalln("repository.NewDatabase: could not open database:", err)
	}

	return Database{
		db: db,
	}
}

func (d Database) Close() error {
	return d.db.Close()
}
