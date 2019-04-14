package server

import (
	"encoding/json"
	"io"
)

type Config struct {
	Host          string          `json:"host"`
	Port          int             `json:"port"`
	StorageType   string          `json:"storageType"`
	StorageConfig json.RawMessage `json:"storage"`
}

func LoadConfig(r io.Reader) (*Config, error) {
	var c Config
	dec := json.NewDecoder(r)
	err := dec.Decode(&c)
	return &c, err
}
