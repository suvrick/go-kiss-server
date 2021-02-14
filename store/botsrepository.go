package store

import "database/sql"

// BotRepository ...
type BotRepository struct {
	db *sql.DB
}

// NewBotRepository ...
func NewBotRepository(db *sql.DB) *BotRepository {
	return &BotRepository{
		db: db,
	}
}
