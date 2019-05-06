package postgres

import (
	"encoding/json"

	"github.com/belljustin/hancock/key"
)

type Config struct {
	key.Config

	User    string `json:"user"`
	Name    string `json:"dbname"`
	SSLMode string `json:"sslmode"`
}

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
