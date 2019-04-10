package fakes

import (
	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
)

type Keys struct{}

func (k Keys) Get(id uuid.UUID) (*models.Key, error) {
	return &models.Key{
		Id:        id,
		Algorithm: "fake",
		Owner:     "belljust.in/justin",
		Pub:       []byte("SGVsbG8sIFdvcmxk"),
		Priv:      []byte{},
	}, nil
}

func (k Keys) Create(owner, alg string, pub, priv []byte) (*models.Key, error) {
	return &models.Key{
		Id:        uuid.New(),
		Algorithm: alg,
		Owner:     owner,
		Pub:       pub,
		Priv:      priv,
	}, nil
}
