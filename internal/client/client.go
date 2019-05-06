package client

import (
	"encoding/json"
	"os"

	"github.com/urfave/cli"

	"github.com/belljustin/hancock/internal/server"
	"github.com/belljustin/hancock/key"
)

type config struct {
	Server  server.Config   `json:"server"`
	Backend string          `json:"backend"`
	Storage json.RawMessage `json:"storage"`
}

func loadConfig(c *cli.Context) (*config, error) {
	f, err := os.Open(c.GlobalString("config"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var conf config
	dec := json.NewDecoder(f)
	if err := dec.Decode(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func getStorage(c config) (key.Storage, error) {
	return key.Open(c.Backend, c.Storage)
}
