package models

import (
	_ "encoding/json"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Keys)
)

func Register(name string, driver Keys) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("keys: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("keys: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Open(name string) (Keys, error) {
	driversMu.RLock()
	defer driversMu.RUnlock()

	s, ok := drivers[name]
	if !ok {
		return nil, errors.New("Key storage not registered")
	}

	err := s.Open()
	return s, err
}

type Key struct {
	Id        uuid.UUID `json:"id"`
	Algorithm string    `json:"alg"`
	Owner     string    `json:"owner"`
	Pub       []byte    `json:"pub"`
	Priv      []byte    `json:"-"`
}

type Keys interface {
	Get(id uuid.UUID) (*Key, error)
	Create(k *Key) error
	Open() error
}
