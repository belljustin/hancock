// Package mem is an in-memory implementation of the `key.Storage` interface
package mem

import (
	"sync"

	"github.com/google/uuid"

	"github.com/belljustin/hancock/key"
)

const (
	driverName = "mem"
)

func init() {
	s := &KeyStorage{}
	key.Register(driverName, s)
}

// KeyStorage implements the key.Storage interface using memory as the storage device. Retrival
// and creation of keys on KeyStorage is the thread-safe.
type KeyStorage struct {
	sync.RWMutex
	m         map[string]key.Key
	generator key.SignerGenerator
}

// Get retrieves a key identified by id from memory.
func (s *KeyStorage) Get(id string) (*key.Key, error) {
	s.RLock()
	defer s.RUnlock()

	k, ok := s.m[id]
	if !ok {
		return nil, nil
	}
	return &k, nil
}

// Create inserts a new key of type alg in memory.
func (s *KeyStorage) Create(alg string, opts key.Opts) (*key.Key, error) {
	s.Lock()
	defer s.Unlock()

	signer, err := s.generator.New(alg, opts)
	if err != nil {
		return nil, err
	}

	k := key.Key{
		Id:        uuid.New().String(),
		Algorithm: alg,
		Signer:    signer,
	}

	s.m[k.Id] = k
	return &k, nil
}

// Open initializes a new in-memory `KeyStorage`. The config is not used and can be left empty.
func (s *KeyStorage) Open(config []byte) error {
	s.m = make(map[string]key.Key)
	s.generator = key.DefaultSignerGenerator
	return nil
}
