package key

import (
	"crypto"
	_ "encoding/json"
)

type Key struct {
	Id        string `json:"id" sql:"id"`
	Algorithm string `json:"alg" sql:"alg"`
	Owner     string `json:"owner" sql:"owner"`
	Signer    crypto.Signer
}

type Opts map[string]interface{}

type KeyStorage interface {
	Get(id string) (*Key, error)
	Create(owner, alg string, o Opts) (*Key, error)
	Open(config []byte) error
}
