package store

import (
	"database/sql"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/model"
)

// UserRepository ...
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository ...
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create ...
func (r *UserRepository) Create(u *model.User) (int, error) {

	if err := u.BeforeCreate(); err != nil {
		return 0, err
	}

	row := r.db.QueryRow("INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id", u.Email, u.EncryptedPassword)

	if row.Err() != nil {
		return 0, row.Err()
	}

	err := row.Scan(&u.ID)
	if err != nil {
		return 0, err
	}

	return u.ID, nil
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {

	u := &model.User{}

	row := r.db.QueryRow("SELECT id, email, encrypted_password FROM users WHERE id = $1", id)
	if row.Err() != nil {
		return nil, errors.ErrRecordNotFound
	}

	err := row.Scan(&u.ID, &u.Email, &u.EncryptedPassword)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	row := r.db.QueryRow("SELECT id, email, encrypted_password FROM users WHERE email = $1", email)

	if row.Err() != nil {
		return nil, errors.ErrRecordNotFound
	}

	err := row.Scan(&u.ID, &u.Email, &u.EncryptedPassword)
	if err != nil {
		return nil, err
	}

	return u, nil
}
