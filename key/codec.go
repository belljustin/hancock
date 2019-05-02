package key

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
)

var DefaultCodec = MultiCodec{
	Codecs: map[string]Codec{
		"rsa": &RsaGobCodec{},
	},
}

type Codec interface {
	Encode(k crypto.Signer) (priv []byte, err error)
	Decode(priv []byte) (s crypto.Signer, err error)
}

type MultiCodec struct {
	Codecs map[string]Codec
}

func (c *MultiCodec) Encode(k crypto.Signer, alg string) ([]byte, error) {
	enc, ok := c.Codecs[alg]
	if !ok {
		return []byte{}, fmt.Errorf("Algorithm '%s' is not supported by the codec.", alg)
	}
	return enc.Encode(k)
}

func (c *MultiCodec) Decode(priv []byte, alg string) (crypto.Signer, error) {
	dec, ok := c.Codecs[alg]
	if !ok {
		return nil, fmt.Errorf("Algorithm '%s' is not supported by the codec.", alg)
	}
	return dec.Decode(priv)
}

// gob

type gobCodec struct{}

func (c *gobCodec) Encode(s crypto.Signer) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(s); err != nil {
		return []byte{}, err
	}
	return b.Bytes(), nil
}

//// rsa

type RsaGobCodec struct {
	*gobCodec
}

func (c *RsaGobCodec) Decode(priv []byte) (crypto.Signer, error) {
	b := bytes.NewBuffer(priv)
	dec := gob.NewDecoder(b)
	var k rsa.PrivateKey
	if err := dec.Decode(&k); err != nil {
		return nil, err
	}
	return &k, nil
}
