package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/belljustin/hancock/server"
)

const (
	url = "http://127.0.0.1:8000"
)

func usage() {
	fmt.Println("hancock usage: TODO")
	os.Exit(-1)
}

// Server Commands

func handleServer(args []string) {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	c, err := server.LoadConfig(f)
	if err != nil {
		log.Fatal(err)
	}

	router := server.NewRouter(c)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.Port), router)
	panic(err)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}
	args := os.Args[1:]

	switch args[0] {
	case "server":
		handleServer(args)
	default:
		usage()
	}
}
