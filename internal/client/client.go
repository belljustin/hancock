package client

import (
	"encoding/json"
	"os"

	"github.com/urfave/cli"

	"github.com/belljustin/hancock/internal/server"
	"github.com/belljustin/hancock/key"
)

type Config struct {
	Server  server.Config   `json:"server"`
	Backend string          `json:"backend"`
	Storage json.RawMessage `json:"storage"`
}

func loadConfig(c *cli.Context) (*Config, error) {
	f, err := os.Open(c.GlobalString("config"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	dec := json.NewDecoder(f)
	if err := dec.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getStorage(c Config) (key.Storage, error) {
	return key.Open(c.Backend, c.Storage)
}
