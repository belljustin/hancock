package key

import (
	"log"
	"os"
)

const (
	// EnvKey is the environment variable specifying secret encryption key.
	EnvKey = "HANCOCK_KEY"

	// RSA the signing algorithm
	RSA = "rsa"

	// AES the encryption algorithm
	AES = "aes"
)

// Config provides configuration for KeyStorage.
type Config struct {
	Encryption string `json:"encryption"`
	Key        string `json:"key"`
}

// LoadEnv replaces empty fields with matching environment variables. See this file's
// constants for a list of available options.
func (c *Config) LoadEnv() {
	if c.Key != "" {
		return
	} else if key, ok := os.LookupEnv(EnvKey); ok {
		c.Key = key
	}
}

// GetCodec returns builtin `MultiCodec`s according to the config's Encryption.
func (c *Config) GetCodec() MultiCodec {
	switch c.Encryption {
	case AES:
		log.Print("hancock: added AES encryption")
		return NewAesCodec(DefaultCodec, c.Key)
	default:
		log.Print("hancock: no encryption enabled")
		return DefaultCodec
	}
}
