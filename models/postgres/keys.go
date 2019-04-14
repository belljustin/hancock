package postgres

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/belljustin/hancock/models"
)

func init() {
	k := &Keys{}
	models.Register("postgres", k)
}

type Keys struct {
	db *sql.DB
}

func (s *Keys) Open() error {
	connStr := "user=hancock dbname=hancock sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	s.db = db
	return s.db.Ping()
}

func (s *Keys) Get(id uuid.UUID) (*models.Key, error) {
	var k models.Key
	query := `SELECT id, alg, owner, pub, priv FROM keys
			  WHERE id = $1`

	r := s.db.QueryRow(query, id)
	if err := r.Scan(&k.Id, &k.Algorithm, &k.Owner, &k.Pub, &k.Priv); err != nil {
		return nil, err
	}
	return &k, nil
}

func (s *Keys) Create(k *models.Key) error {
	update := `INSERT INTO keys(id, alg, owner, pub, priv)
			   VALUES($1, $2, $3, $4, $5)`

	res, err := s.db.Exec(update, k.Id, k.Algorithm, k.Owner, k.Pub, k.Priv)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n <= 0 {
		return errors.New("No rows updated")
	}
	return nil
}
