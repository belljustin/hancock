package mem

import (
	"sync"

	"github.com/google/uuid"

	"github.com/belljustin/hancock/models"
)

func init() {
	k := Keys{}
	models.Register("mem", &k)
}

type Keys struct {
	sync.RWMutex
	m map[uuid.UUID]models.Key
}

func (s *Keys) Get(id uuid.UUID) (*models.Key, error) {
	s.RLock()
	defer s.RUnlock()

	k, ok := s.m[id]
	if !ok {
		return nil, nil
	}
	return &k, nil
}

func (s *Keys) Create(k *models.Key) error {
	s.Lock()
	defer s.Unlock()

	s.m[k.Id] = *k
	return nil
}

func (s *Keys) Open(config []byte) error {
	s.m = make(map[uuid.UUID]models.Key)
	return nil
}
