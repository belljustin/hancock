package key

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"io"
)

var DefaultCodec MultiCodec

func init() {
	DefaultCodec = &multiCodec{
		map[string]Codec{
			"rsa": &RsaGobCodec{},
		},
	}
}

type Codec interface {
	Encode(k crypto.Signer) (priv []byte, err error)
	Decode(priv []byte) (s crypto.Signer, err error)
}

type MultiCodec interface {
	Encode(k crypto.Signer, alg string) (priv []byte, err error)
	Decode(priv []byte, alg string) (s crypto.Signer, err error)
}

type multiCodec struct {
	Codecs map[string]Codec
}

func (c *multiCodec) Encode(k crypto.Signer, alg string) ([]byte, error) {
	enc, ok := c.Codecs[alg]
	if !ok {
		return []byte{}, fmt.Errorf("Algorithm '%s' is not supported by the codec.", alg)
	}
	return enc.Encode(k)
}

func (c *multiCodec) Decode(priv []byte, alg string) (crypto.Signer, error) {
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

// encryption

//// AES

type AesCodec struct {
	key                 []byte
	clearTextMultiCodec MultiCodec
}

func NewAesCodec(clearTextMultiCodec MultiCodec, key string) *AesCodec {
	hasher := md5.New()
	hasher.Write([]byte(key))
	hKey := hasher.Sum(nil)

	return &AesCodec{
		key:                 hKey,
		clearTextMultiCodec: clearTextMultiCodec,
	}
}

func (c *AesCodec) Encode(s crypto.Signer, alg string) ([]byte, error) {
	priv, err := c.clearTextMultiCodec.Encode(s, alg)
	if err != nil {
		return []byte{}, err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return []byte{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}
	ciphertext := gcm.Seal(nonce, nonce, priv, nil)
	return ciphertext, nil
}

func (c *AesCodec) Decode(data []byte, alg string) (crypto.Signer, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	priv, err := gcm.Open(nil, nonce, ciphertext, nil)

	return c.clearTextMultiCodec.Decode(priv, alg)
}
