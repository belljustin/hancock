package key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
)

// GenerateSignerFunc is a type of function that produces a Signer given some `Opts`.
type GenerateSignerFunc func(o Opts) (crypto.Signer, error)

// SignerGenerator collects `GenerateSignerFunc`s by the cryptographic signing algorithms on which
// they rely.
type SignerGenerator struct {
	Generators map[string]GenerateSignerFunc
}

// New generates a new Signer using the cryptographic signing algorithm specified by alg using
// the provided `Opts`.
func (f *SignerGenerator) New(alg string, o Opts) (crypto.Signer, error) {
	g, ok := f.Generators[alg]
	if !ok {
		return nil, fmt.Errorf("algorithm '%s' is not a supported signer alg", alg)
	}
	return g(o)
}

// DefaultSignerGenerator is a sensible default generator that supports RSA.
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
