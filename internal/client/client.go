package client

import (
	"os"

	"github.com/urfave/cli"

	"github.com/belljustin/hancock/internal/server"
	"github.com/belljustin/hancock/key"
)

func loadConfig(c *cli.Context) (*server.Config, error) {
	f, err := os.Open(c.GlobalString("config"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return server.LoadConfig(f)
}

func getStorage(c server.Config) (key.Storage, error) {
	return key.Open(c.StorageType, c.StorageConfig)
}
