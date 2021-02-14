package store

import (
	"database/sql"
	"github.com/suvrick/go-kiss-server/errors"
	"github.com/suvrick/go-kiss-server/model"
)

// ProxyRepository ...
type ProxyRepository struct {
	db *sql.DB
}

// NewProxyRepository ...
func NewProxyRepository(db *sql.DB) *ProxyRepository {
	return &ProxyRepository{
		db: db,
	}
}

// Create ...
func (r *ProxyRepository) Create(p *model.Proxy) (int, error) {
	row := r.db.QueryRow(
		"INSERT INTO proxies (host, login, password, isbad, isbusy, dataadd, datalastuse) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		&p.Host,
		&p.Login,
		&p.Password,
		&p.IsBad,
		&p.IsBusy,
		&p.DataAdd,
		&p.DataLastUse,
	)

	if row.Err() != nil {
		return 0, row.Err()
	}

	err := row.Scan(&p.ID)
	if err != nil {
		return 0, err
	}

	return p.ID, nil
}

// Find ...
func (r *ProxyRepository) Find(proxyID int) (*model.Proxy, error) {
	p := &model.Proxy{}
	row := r.db.QueryRow("SELECT id, host, login, password, isbad, isbusy, dataadd, datalastuse FROM proxies WHERE id = $1", proxyID)

	if row.Err() != nil {
		return nil, errors.ErrRecordNotFound
	}

	err := row.Scan(
		&p.ID,
		&p.Host,
		&p.Login,
		&p.Password,
		&p.IsBad,
		&p.IsBusy,
		&p.DataAdd,
		&p.DataLastUse,
	)

	if err != nil {
		return nil, errors.ErrRecordNotFound
	}

	return p, nil
}

// Free ...
func (r *ProxyRepository) Free() (*model.Proxy, error) {
	p := &model.Proxy{}
	row := r.db.QueryRow("SELECT id, host, login, password, isbad, isbusy, dataadd, datalastuse FROM proxies WHERE isbad != true and isbusy != true and count != 2")

	if row.Err() != nil {
		return nil, errors.ErrRecordNotFound
	}

	err := row.Scan(
		&p.ID,
		&p.Host,
		&p.Login,
		&p.Password,
		&p.IsBad,
		&p.IsBusy,
		&p.DataAdd,
		&p.DataLastUse,
	)

	if err != nil {
		return nil, errors.ErrRecordNotFound
	}

	return p, nil
}

// All ...
func (r *ProxyRepository) All() ([]*model.Proxy, error) {

	ps := make([]*model.Proxy, 0)
	rows, err := r.db.Query("SELECT id, host, login, password, isbad, isbusy, dataadd, datalastuse FROM proxies")

	if err != nil {
		return nil, errors.ErrRecordNotFound
	}

	defer rows.Close()

	for rows.Next() {
		p := &model.Proxy{}
		err := rows.Scan(
			&p.ID,
			&p.Host,
			&p.Login,
			&p.Password,
			&p.IsBad,
			&p.IsBusy,
			&p.DataAdd,
			&p.DataLastUse,
		)

		if err != nil {
			return nil, errors.ErrRecordNotFound
		}

		ps = append(ps, p)
	}

	return ps, nil
}
