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

type KeyStorage struct {
	sync.RWMutex
	m         map[string]key.Key
	generator key.SignerGenerator
}

func (s *KeyStorage) Get(id string) (*key.Key, error) {
	s.RLock()
	defer s.RUnlock()

	k, ok := s.m[id]
	if !ok {
		return nil, nil
	}
	return &k, nil
}

func (s *KeyStorage) Create(owner, alg string, opts key.Opts) (*key.Key, error) {
	s.Lock()
	defer s.Unlock()

	signer, err := s.generator.New(alg, opts)
	if err != nil {
		return nil, err
	}

	k := key.Key{
		Id:        uuid.New().String(),
		Algorithm: alg,
		Owner:     owner,
		Signer:    signer,
	}

	s.m[k.Id] = k
	return &k, nil
}

func (s *KeyStorage) Open(config []byte) error {
	s.m = make(map[string]key.Key)
	s.generator = key.DefaultSignerGenerator
	return nil
}
