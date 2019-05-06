// Package postgres is postgres database implementation of the `key.Storage` interface
package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/belljustin/hancock/key"
)

func init() {
	s := &KeyStorage{}
	key.Register("postgres", s)
}

// KeyStorage is an implementation of `key.Storage` using a postgres as a backend.
type KeyStorage struct {
	db *sql.DB

	codec     key.MultiCodec
	generator key.SignerGenerator
}

// Open configures the `KeyStorage` using rawConfig and connects to the database.
// If validation passes, the return value is the result of pinging the database.
func (s *KeyStorage) Open(rawConfig []byte) error {
	c, err := LoadConfig(rawConfig)
	if err != nil {
		return err
	}

	s.codec = c.GetCodec()
	s.generator = key.DefaultSignerGenerator

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", c.User, c.Name, c.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	s.db = db
	return s.db.Ping()
}

// Get fetches the `key.Key` specified by the unique sid from the database. All IDs MUST parse to
// a valid uuid.
func (s *KeyStorage) Get(sid string) (*key.Key, error) {
	var k key.Key
	query := `SELECT id, alg, owner, priv FROM keys
			  WHERE id = $1`

	id, err := uuid.Parse(sid)
	if err != nil {
		return nil, err
	}

	var data []byte
	r := s.db.QueryRow(query, id)
	if err := r.Scan(&k.Id, &k.Algorithm, &k.Owner, &data); err != nil {
		return nil, err
	}

	signer, err := s.codec.Decode(data, k.Algorithm)
	if err != nil {
		return nil, err
	}
	k.Signer = signer

	return &k, nil
}

// Create inserts a new `key.Key` into the database. The id will be generated as a v4 uuid.
func (s *KeyStorage) Create(owner, alg string, opts key.Opts) (*key.Key, error) {
	update := `INSERT INTO keys(id, alg, owner, priv)
			   VALUES($1, $2, $3, $4)`

	signer, err := s.generator.New(alg, opts)
	if err != nil {
		return nil, err
	}

	data, err := s.codec.Encode(signer, alg)
	if err != nil {
		return nil, err
	}

	id := uuid.New()

	res, err := s.db.Exec(update, id, alg, owner, data)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if n <= 0 {
		return nil, errors.New("No rows updated")
	}

	return &key.Key{
		Id:        id.String(),
		Algorithm: alg,
		Owner:     owner,
		Signer:    signer,
	}, nil
}
