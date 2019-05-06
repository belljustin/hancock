// Package key is an interface for opaque cryptographic signing keys
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

// Register makes a key `Storage` available by the provided name. If Register is called twice
// with the same name or if driver is nil, it panics.
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

// Open opens a key `Storage` specified by its storage name and configuration. The configuration
// is usually a []byte representation of a json config.
//
// Most users will open storage via a driver-specific connetion helper function that returns a Storage.
// Hancock includes a few drivers which are subpackages of key.
//
// Open may just validate its arguments without creating a connection to the Storage.
func Open(driverName string, config []byte) (Storage, error) {
	driversMu.RLock()
	driver, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("hancock: unknown driver %s (forgotten import?)", driverName)
	}
	return driver, driver.Open(config)
}

// Key is an interface for opaque cryptographic signing keys.
type Key struct {
	// Id is a unique identifier for the key.
	Id string `json:"id" sql:"id"`
	// Algorithm specifies the cryptographic signing algorithm underlying this key.
	Algorithm string `json:"alg" sql:"alg"`
	// Owner identifies the owner of this key.
	Owner string `json:"owner" sql:"owner"`
	// Signer implements the crypto.Signer interface which can be used for signing and inspecting
	// the public key.
	Signer crypto.Signer
}

// Opts specify additional options used in `Key` generation.
type Opts map[string]interface{}

// Storage is an interface for a storage backend of `Key`s. Some implementations of Storage can
// be found as subpackages of key.
type Storage interface {
	// Get retrieves a `*Key` using the unique identifier id. If no key with that id is found, both
	// return values are null.
	Get(id string) (*Key, error)
	// Create inserts a new `Key` generated using the algorithm specified by alg and the provided
	// `Opts`. The resulting `Key` is returned.
	Create(owner, alg string, o Opts) (*Key, error)
	// Open opens a key storage. This must be called before calling other methods on `Storage`.
	// Most users will Open a key `Storage` using the a driverName as in `Open`.
	Open(config []byte) error
}
