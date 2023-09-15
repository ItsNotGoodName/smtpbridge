package database

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

type Querier interface {
	Conn() *sql.DB
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	BeginTx(ctx context.Context, write bool) (QuerierTx, error)
}

type QuerierTx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Commit() error
	Rollback() error
}

func New(dbPath string, debug bool) (Querier, error) {
	// https://github.com/pocketbase/pocketbase/blob/2608efb56cd4dab3cccfea5d6822cb87f7308636/core/db_nocgo.go#L14
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-16000)"
	db, err := sql.Open("sqlite", dbPath+pragmas)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := dummyCreate(db); err != nil {
		return nil, err
	}

	if debug {
		return DebugDB{DB: db}, nil
	}

	return DB{DB: db}, nil
}

type DB struct {
	*sql.DB
}

func (db DB) BeginTx(ctx context.Context, write bool) (QuerierTx, error) {
	return dummyBeginTx(ctx, db.DB, write)
}

func (db DB) Conn() *sql.DB {
	return db.DB
}

type DebugDB struct {
	*sql.DB
}

func NewDebugDB(db *sql.DB) DebugDB {
	return DebugDB{db}
}

func (db DebugDB) Conn() *sql.DB {
	return db.DB
}

func (db DebugDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Debug().
		Str("func", "ExecContext").
		Any("args", args).
		Msg(query)
	return db.DB.ExecContext(ctx, query, args...)
}

func (tx DebugTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Debug().
		Str("func", "ExecContext (Tx)").
		Any("args", args).
		Msg(query)
	return tx.Tx.ExecContext(ctx, query, args...)
}

func (db DebugDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Debug().
		Str("func", "QueryContext").
		Any("args", args).
		Msg(query)
	return db.DB.QueryContext(ctx, query, args...)
}

func (tx DebugTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Debug().
		Str("func", "QueryContext (Tx)").
		Any("args", args).
		Msg(query)
	return tx.Tx.QueryContext(ctx, query, args...)
}

func (db DebugDB) BeginTx(ctx context.Context, write bool) (QuerierTx, error) {
	log.Debug().
		Msg("BeginTx (Tx)")
	tx, err := dummyBeginTx(ctx, db.DB, write)
	if err != nil {
		return DebugTx{}, err
	}
	return DebugTx{Tx: tx}, nil
}

func (tx DebugTx) Commit() error {
	log.Debug().
		Str("func", "Commit (Tx)").
		Msg("")
	return tx.Tx.Commit()
}

func (tx DebugTx) Rollback() error {
	log.Debug().
		Str("func", "Rollback (Tx)").
		Msg("")
	return tx.Tx.Rollback()
}

type DebugTx struct {
	*sql.Tx
}
