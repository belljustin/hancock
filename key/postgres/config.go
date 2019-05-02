package postgres

import (
	"encoding/json"
)

type Config struct {
	User    string `json:"user"`
	Name    string `json:"dbname"`
	SSLMode string `json:"sslmode"`
}

func LoadConfig(rawConfig []byte) (*Config, error) {
	var c Config
	err := json.Unmarshal(rawConfig, &c)
	return &c, err
}
