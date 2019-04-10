package mem

import (
	"sync"

	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
)

func init() {
	models.Register("mem", Keys{})
}

type Keys struct {
	sync.RWMutex
	m map[uuid.UUID]models.Key
}

func (s Keys) Get(id uuid.UUID) (*models.Key, error) {
	s.RLock()
	defer s.RUnlock()

	k, ok := s.m[id]
	if !ok {
		return nil, nil
	}
	return &k, nil
}

func (s Keys) Create(alg, owner string, pub, priv []byte) (*models.Key, error) {
	s.Lock()
	defer s.Unlock()

	k := models.Key{
		Id:        uuid.New(),
		Algorithm: alg,
		Owner:     owner,
		Pub:       pub,
		Priv:      priv,
	}
	s.m[k.Id] = k

	return &k, nil
}
