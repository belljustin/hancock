package models

import (
	_ "encoding/json"

	"github.com/google/uuid"
)

type Key struct {
	Id        uuid.UUID `json:"id"`
	Algorithm string    `json:"alg"`
	Owner     string    `json:"owner"`
	Pub       []byte    `json:"pub"`
	Priv      []byte    `json:"-"`
}

type Keys interface {
	Get(id uuid.UUID) (*Key, error)
	Create(owner string) (*Key, error)
}
