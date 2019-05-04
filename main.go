package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/belljustin/hancock/internal/client"
)

const (
	url = "http://127.0.0.1:8000"
)

func main() {
	app := cli.NewApp()
	app.Name = "hancock"
	app.Usage = "Create and use cryptographic signing keys"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "./config.json",
			Usage: "path to config file",
		},
	}

	app.Commands = []cli.Command{
		client.ServerCmd,
		client.KeyCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
