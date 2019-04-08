package fakes

import (
	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
)

type Keys struct{}

func (k Keys) Get(id uuid.UUID) (*models.Key, error) {
	return &models.Key{
		Id:    id,
		Owner: "belljust.in/justin",
		Data:  []byte{},
	}, nil
}

func (k Keys) Create(owner string) (*models.Key, error) {
	return &models.Key{
		Id:    uuid.New(),
		Owner: owner,
		Data:  []byte{},
	}, nil
}
