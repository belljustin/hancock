package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
)

type Alg interface {
	NewKey(owner string) (*models.Key, error)
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
