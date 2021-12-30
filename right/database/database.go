package database

import (
	"io/fs"
	"log"
	"os"

	"github.com/asdine/storm"
)

type DB struct {
	db     *storm.DB
	attDir string
	fs     fs.FS
}

func NewDB(dbFile, attDir string) *DB {
	db, err := storm.Open(dbFile)
	if err != nil {
		log.Fatal("database.NewDB:", err)
	}

	err = os.MkdirAll(attDir, 0755)
	if err != nil {
		log.Fatal("database.NewDB:", err)
	}

	return &DB{
		db:     db,
		fs:     os.DirFS(attDir),
		attDir: attDir,
	}
}
