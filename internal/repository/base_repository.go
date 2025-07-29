package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db *sqlx.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sqlx.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB returns the database connection
func (r *BaseRepository) GetDB() *sqlx.DB {
	return r.db
}

// WithTransaction executes a function within a database transaction
func (r *BaseRepository) WithTransaction(fn func(*sqlx.Tx) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback() // Ignore rollback errors in defer
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
