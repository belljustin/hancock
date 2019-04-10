package fakes

import (
	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
)

type Keys struct{}

func (s Keys) Get(id uuid.UUID) (*models.Key, error) {
	return &models.Key{
		Id:        id,
		Algorithm: "fake",
		Owner:     "belljust.in/justin",
		Pub:       []byte("SGVsbG8sIFdvcmxk"),
		Priv:      []byte{},
	}, nil
}

func (s Keys) Create(k *models.Key) error {
	return nil
}
