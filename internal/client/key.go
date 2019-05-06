package client

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/urfave/cli"

	"github.com/belljustin/hancock/key"
)

var hashes = map[string]crypto.Hash{
	"sha256": crypto.SHA256,
}

var KeyCmd = cli.Command{
	Name:  "key",
	Usage: "manage keys",
	Subcommands: []cli.Command{
		createKeyCmd,
		getKeyCmd,
		signCmd,
	},
}

type clientFunc func(key.Storage, *cli.Context) error

func createClientFunc(f clientFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		conf, err := loadConfig(c)
		if err != nil {
			return err
		}

		storage, err := key.Open(conf.Backend, conf.Storage)
		if err != nil {
			return err
		}

		return f(storage, c)
	}
}

var createKeyCmd = cli.Command{
	Name:   "create",
	Usage:  "create a new key",
	Action: createClientFunc(createKey),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "alg",
			Usage: "the algorithm to use in key generation",
		},
	},
}

func createKey(s key.Storage, c *cli.Context) error {
	alg := c.String("alg")
	if alg == "" {
		return errors.New("alg must not be empty")
	}

	opts := key.Opts{}
	k, err := s.Create("belljust.in/justin", alg, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Created key %s", k.Id)
	return nil
}

var getKeyCmd = cli.Command{
	Name:   "get",
	Usage:  "fetch an existing key",
	Action: createClientFunc(getKey),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "the key identifier",
		},
	},
}

func getKey(s key.Storage, c *cli.Context) error {
	id := c.String("id")
	if id == "" {
		return errors.New("id must not be empty")
	}

	k, err := s.Get(id)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", k.Signer.Public())
	return nil
}

var signCmd = cli.Command{
	Name:   "sign",
	Usage:  "sign a digest",
	Action: createClientFunc(sign),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "the key identifier",
		},
		cli.StringFlag{
			Name:  "digest",
			Usage: "the hash of data to be signed",
		},
		cli.StringFlag{
			Name:  "hash",
			Usage: "the hashing algorithm used to create the digest",
			Value: "sha256",
		},
	},
}

func sign(s key.Storage, c *cli.Context) error {
	id := c.String("id")
	if id == "" {
		return errors.New("id must not be empty")
	}

	digest := c.String("digest")
	if digest == "" {
		return errors.New("digest must not be empty")
	}
	bDigest, err := hex.DecodeString(digest)
	if err != nil {
		return err
	}

	h := c.String("hash")
	hash, ok := hashes[h]
	if !ok {
		return fmt.Errorf("Hash '%s' is not supported", h)
	}

	k, err := s.Get(id)
	if err != nil {
		return err
	}

	signature, err := k.Signer.Sign(rand.Reader, bDigest, hash)
	if err != nil {
		return err
	}

	fmt.Printf("%v", signature)
	return nil
}
