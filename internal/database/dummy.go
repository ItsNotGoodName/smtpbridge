package database

import (
	"context"
	"database/sql"
)

func dummyCreate(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS _dummy(_dummy integer);")
	return err
}

func dummyBeginTx(ctx context.Context, db *sql.DB, write bool) (*sql.Tx, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	if write {
		// This prevents SQLITE_BUSY (5) and database locked (517) when doing write transactions
		// because we tell sqlite that we are going to do a write transaction through the dummy DELETE query.
		_, err = tx.ExecContext(ctx, "DELETE FROM _dummy WHERE 0 = 1;")
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}
