package key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
)

type GenerateSignerFunc func(o Opts) (crypto.Signer, error)

type SignerGenerator struct {
	Generators map[string]GenerateSignerFunc
}

func (f *SignerGenerator) New(alg string, o Opts) (crypto.Signer, error) {
	g, ok := f.Generators[alg]
	if !ok {
		return nil, fmt.Errorf("Algorithm '%s' is not a supported signer alg.", alg)
	}
	return g(o)
}

var DefaultSignerGenerator = SignerGenerator{
	Generators: map[string]GenerateSignerFunc{
		RSA: rsaGenerateSigner,
	},
}

// RSA

func rsaGenerateSigner(o Opts) (crypto.Signer, error) {
	var bits int
	b, ok := o["bits"]
	if !ok {
		bits = 2048
	} else {
		if bits, ok = b.(int); !ok {
			return nil, errors.New("Could not cast bits to int")
		}
	}

	return rsa.GenerateKey(rand.Reader, bits)
}
