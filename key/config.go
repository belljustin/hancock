package key

import (
	"os"
)

const (
	// environment variables
	ENV_KEY = "HANCOCK_KEY"

	// signing algorithms
	RSA = "rsa"
)

type Config struct {
	Encryption string `json:"encryption"`
	Key        string `json:"key"`
}

func (c *Config) LoadEnv() {
	if c.Key != "" {
		return
	} else if key, ok := os.LookupEnv(ENV_KEY); ok {
		c.Key = key
	}
}

func (c *Config) GetCodec() MultiCodec {
	switch c.Encryption {
	case RSA:
		return NewAesCodec(DefaultCodec, c.Key)
	default:
		return DefaultCodec
	}
}
