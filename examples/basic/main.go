package main

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/belljustin/hancock/key"
	_ "github.com/belljustin/hancock/key/mem" // In memory driver
)

func main() {
	s, _ := key.Open("mem", []byte{})

	k, _ := s.Create("owner", "rsa", nil)

	hasher := sha256.New()
	hasher.Write([]byte("document to sign"))
	digest := hasher.Sum(nil)

	signature, _ := k.Signer.Sign(rand.Reader, digest, crypto.SHA256)
	fmt.Printf("%x", signature)
}
