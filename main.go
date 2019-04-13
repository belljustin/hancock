package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"

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
	fmt.Println("hancock key usage: ")
	os.Exit(-1)
}

func handleKeys(args []string) {
	if len(args) == 1 {
		keyUsage()
	}
	args = args[1:]

	switch args[0] {
	case "new":
		handleNewKey(args)
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
		fmt.Println(err)
		os.Exit(-1)
	}

	if alg == "" {
		fmt.Println("Required flag alg cannot be empty")
		os.Exit(-1)
	}

	c := client.NewHancockClient(url)
	k, err := c.NewKey(alg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	pubKey := base64.StdEncoding.EncodeToString(k.Pub)
	fmt.Printf("id: %s\nalg: %s\nowner: %s\npub: %s\n", k.Id.String(), k.Algorithm, k.Owner, pubKey)
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