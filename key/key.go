package key

import (
	"crypto"
	_ "encoding/json"
	"fmt"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Storage)
)

func Register(name string, driver Storage) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("hancock: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("hancock: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Open(driverName string, config []byte) (Storage, error) {
	driversMu.RLock()
	driver, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("hancock: unknown driver %s (forgotten import?)", driverName)
	}
	return driver, driver.Open(config)
}

type Key struct {
	Id        string `json:"id" sql:"id"`
	Algorithm string `json:"alg" sql:"alg"`
	Owner     string `json:"owner" sql:"owner"`
	Signer    crypto.Signer
}

type Opts map[string]interface{}

type Storage interface {
	Get(id string) (*Key, error)
	Create(owner, alg string, o Opts) (*Key, error)
	Open(config []byte) error
}
