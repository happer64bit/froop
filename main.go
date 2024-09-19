package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"

	"froop/server"
)

func main() {
	parser := argparse.NewParser("myserver", "A simple HTTP server")

	serveCmd := parser.NewCommand("serve", "Start the HTTP server")

	port := serveCmd.String("p", "port", &argparse.Options{
		Help:    "Port to run the server on",
		Default: "8080",
	})
	address := serveCmd.String("a", "address", &argparse.Options{
		Help:    "Address to bind the server to",
		Default: "localhost",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	if len(os.Args) < 2 || os.Args[1] != "serve" {
		fmt.Println("Usage: <program> serve [--port=PORT] [--address=ADDRESS]")
		os.Exit(1)
	}

	server.StartServer(*address, *port, ".")
}
