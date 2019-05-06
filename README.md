CAUTION: still in active development and subject to change.

# Hancock

An interface for opaque cryptographic signatures.

Use the standard library `crypto.Signer` interface for signing with keys stored in-memory, a database, or an HSM managed by your cloud provider.
Hancock also provides a binary that includes a CLI and a json REST server.

## Currently Supported Drivers

- mem
- postgres

## Basic Usage

```go
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
	// Create a digest of the document to be signed
	hasher := sha256.New()
	hasher.Write([]byte("document to be signed"))
	digest := hasher.Sum(nil)

	// Open a key storage
	s, _ := key.Open("mem", []byte{})

	// Create a new RSA signing key
	k, _ := s.Create("owner", "rsa", nil)

	// Use k.Signer implements `crypto.Signer`
	signature, _ := k.Signer.Sign(rand.Reader, digest, crypto.SHA256)
	fmt.Printf("%x", signature)
}
```

For more info, please reference the godoc.

## CLI

Run hancock without arguments for usage information.

```
NAME:
   hancock - Create and use cryptographic signing keys

USAGE:
   hancock [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     server, s  start a hancock REST server
     key        manage keys
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value  path to config file (default: "./config.json")
   --help, -h      show help
   --version, -v   print the version
```

### Server	

The hancock server exposes the `key.Storage` interface as a json REST server.

### Key

The hancock key command exposes the `key.Storage` interface as a CLI.
