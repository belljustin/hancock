package client

import (
	"github.com/urfave/cli"

	"github.com/belljustin/hancock/internal/server"
	"github.com/belljustin/hancock/key"
	_ "github.com/belljustin/hancock/key/mem"
	_ "github.com/belljustin/hancock/key/postgres"
)

var ServerCmd = cli.Command{
	Name:    "server",
	Aliases: []string{"s"},
	Usage:   "start a hancock REST server",
	Action:  runServer,
}

func runServer(c *cli.Context) error {
	conf, err := loadConfig(c)
	if err != nil {
		return err
	}

	s, err := key.Open(conf.Backend, conf.Storage)
	if err != nil {
		return err
	}
	return server.Run(conf.Server.Port, s)
}
