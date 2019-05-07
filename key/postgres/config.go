package postgres

import (
	"encoding/json"
	"os"

	"github.com/belljustin/hancock/key"
)

const (
	// EnvPostgresPassword is the environment variable name for the postgres password.
	EnvPostgresPassword = "HANCOCK_POSTGRES_PASSWORD"
)

// Config is a struct for holding settings for a postgres backed key `Storage`.
type Config struct {
	key.Config

	// The database user
	User string `json:"user"`
	// The user password
	Password string `json:"password"`

	// The hostname of the database
	Host string `json:"host"`
	// The port of the database
	Port int `json:"port"`
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

	// If host is the zero value, use localhost
	if c.Host == "" {
		c.Host = "localhost"
	}

	// If port is the zero value, use default postgres port
	if c.Port == 0 {
		c.Port = 5432
	}

	c.loadEnv()
	return &c, err
}

func (c *Config) loadEnv() {
	c.LoadEnv()

	if key, ok := os.LookupEnv(EnvPostgresPassword); c.Password == "" && ok {
		c.Key = key
	}
}
