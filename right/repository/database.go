package repository

import (
	"log"

	"github.com/ItsNotGoodName/smtpbridge/config"
	"github.com/asdine/storm"
)

func NewStorm(cfg *config.Config) *storm.DB {
	db, err := storm.Open(cfg.DB.DB)
	if err != nil {
		log.Fatalln("repository.NewStorm: could not open database:", err)
	}
	return db
}
