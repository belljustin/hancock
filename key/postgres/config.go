package postgres

import (
	"encoding/json"

	"github.com/belljustin/hancock/key"
)

// Config is a struct for holding settings for a postgres backed key `Storage`.
type Config struct {
	key.Config

	// The database user
	User string `json:"user"`
	// The database name
	Name string `json:"dbname"`
	// The ssl mode for connecting to the database.
	// See https://www.postgresql.org/docs/9.1/ssl-tcp.html for options and info.
	SSLMode string `json:"sslmode"`
}

// LoadConfig loads the config provided in the []byte rawConfig. It is assummed the array
// encodes a json configuration of `Config`.
func LoadConfig(rawConfig []byte) (*Config, error) {
	var c Config
	err := json.Unmarshal(rawConfig, &c)
	if err != nil {
		return nil, err
	}
	c.loadEnv()
	return &c, err
}

func (c *Config) loadEnv() {
	c.LoadEnv()
}
