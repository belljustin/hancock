package models

import (
	_ "encoding/json"

	"github.com/google/uuid"
)

type Key struct {
	Id    uuid.UUID `json:"id"`
	Owner string    `json:"owner"`
	Data  []byte    `json:"-"`
}

type Keys interface {
	Get(id uuid.UUID) (*Key, error)
	Create(owner string) (*Key, error)
}
