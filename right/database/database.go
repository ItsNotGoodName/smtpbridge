package database

import (
	"io/fs"
	"log"
	"os"

	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/asdine/storm"
)

type DB struct {
	db     *storm.DB
	attDir string
	fs     fs.FS
}

func New(config *domain.ConfigDB) *DB {
	db, err := storm.Open(config.DB)
	if err != nil {
		log.Fatalln("database.New: could not open database:", err)
	}

	err = os.MkdirAll(config.Attachments, 0755)
	if err != nil {
		log.Fatalln("database.New: could not create attachments directory:", err)
	}

	return &DB{
		db:     db,
		fs:     os.DirFS(config.Attachments),
		attDir: config.Attachments,
	}
}
