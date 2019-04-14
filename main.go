package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/belljustin/hancock/client"
	"github.com/belljustin/hancock/server"
)

const (
	url = "http://127.0.0.1:8000"
)

func usage() {
	fmt.Println("hancock usage: ")
	os.Exit(-1)
}

// Server Commands

func handleServer(args []string) {
	router := server.NewRouter()
	err := http.ListenAndServe(":8000", router)
	panic(err)
}

// Key Commands

func keyUsage() {
	log.Fatal("hancock key usage: ")
}

func handleKeys(args []string) {
	if len(args) == 1 {
		keyUsage()
	}
	args = args[1:]

	switch args[0] {
	case "new":
		handleNewKey(args)
	case "get":
		handleGetKey(args)
	case "sign":
		handleCreateSignature(args)
	default:
		keyUsage()
	}
}

func handleNewKey(args []string) {
	var alg string

	cmd := flag.NewFlagSet("newKey", flag.ExitOnError)
	cmd.StringVar(&alg, "alg", "", "Algorithm used to create the new key")
	err := cmd.Parse(args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if alg == "" {
		log.Fatal("Required flag alg cannot be empty")
	}

	c := client.NewHancockClient(url)
	k, err := c.NewKey(alg)
	if err != nil {
		log.Fatal(err)
	}

	pubKey := base64.StdEncoding.EncodeToString(k.Pub)
	fmt.Printf("id: %s\nalg: %s\nowner: %s\npub: %s\n", k.Id.String(), k.Algorithm, k.Owner, pubKey)
}

func handleGetKey(args []string) {
	var sid string

	cmd := flag.NewFlagSet("getKey", flag.ExitOnError)
	cmd.StringVar(&sid, "id", "", "Key identifier")
	err := cmd.Parse(args[1:])
	if err != nil {
		log.Fatal(err)
	}

	id, err := uuid.Parse(sid)
	if err != nil {
		log.Fatal(err)
	}

	c := client.NewHancockClient(url)
	k, err := c.GetKey(id)
	if err != nil {
		log.Fatal(err)
	}

	pubKey := base64.StdEncoding.EncodeToString(k.Pub)
	fmt.Printf("id: %s\nalg: %s\nowner: %s\npub: %s\n", k.Id.String(), k.Algorithm, k.Owner, pubKey)
}

func handleCreateSignature(args []string) {
	var sid string
	var hash string
	var digest string

	cmd := flag.NewFlagSet("createSignature", flag.ExitOnError)
	cmd.StringVar(&sid, "id", "", "Key identifier")
	cmd.StringVar(&hash, "hash", "sha256", "Hashing algorithm used to produce digest")
	cmd.StringVar(&digest, "digest", "", "Digest to be signed")
	err := cmd.Parse(args[1:])
	if err != nil {
		log.Fatal(err)
	}

	id, err := uuid.Parse(sid)
	if err != nil {
		log.Fatal(err)
	}

	if digest == "" {
		log.Fatal("digest must not be empty")
	}

	c := client.NewHancockClient(url)
	s, err := c.CreateSignature(id, hash, digest)
	if err != nil {
		log.Fatal(err)
	}

	b64s := base64.StdEncoding.EncodeToString(s)
	fmt.Printf("signature: %s\n", b64s)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}
	args := os.Args[1:]

	switch args[0] {
	case "server":
		handleServer(args)
	case "keys":
		handleKeys(args)
	default:
		usage()
	}
}
