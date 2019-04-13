package server

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"

	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
)

var hashes = map[string]crypto.Hash{
	"sha256": crypto.SHA256,
}

type Alg interface {
	NewKey(owner string) (*models.Key, error)
	Sign(priv, digest []byte, hashName string) ([]byte, error)
}

type Rsa struct{}

func (a *Rsa) NewKey(owner string) (*models.Key, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	priv := x509.MarshalPKCS1PrivateKey(key)
	pub := x509.MarshalPKCS1PublicKey(&key.PublicKey)

	return &models.Key{
		Id:        uuid.New(),
		Algorithm: "rsa",
		Owner:     owner,
		Pub:       pub,
		Priv:      priv,
	}, nil
}

func (a *Rsa) Sign(priv, digest []byte, hashName string) ([]byte, error) {
	h, ok := hashes[hashName]
	if !ok {
		return nil, errors.New("Invalid hash")
	}

	k, err := x509.ParsePKCS1PrivateKey(priv)
	if err != nil {
		return nil, err
	}

	s, err := k.Sign(rand.Reader, digest, h)
	return s, err
}
