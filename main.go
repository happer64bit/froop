package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"

	"github.com/happer64bit/froop/server"
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

	auth := serveCmd.String("", "auth", &argparse.Options{
		Help:     "Enable authentication with the format username:password",
		Required: false,
	})

	verbose := serveCmd.Flag("v", "verbose", &argparse.Options{
		Help: "Enable verbose logging",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	if len(os.Args) < 2 || os.Args[1] != "serve" {
		fmt.Println("Usage: <program> serve [--port=PORT] [--address=ADDRESS] [--auth=username:password] [--verbose]")
		os.Exit(1)
	}

	var username, password string
	if *auth != "" {
		parts := strings.SplitN(*auth, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Error: --auth should be in the format username:password.")
			fmt.Println(parser.Usage(nil))
			os.Exit(1)
		}
		username = parts[0]
		password = parts[1]
	}

	// Start the server with or without authentication based on the presence of the --auth flag
	server.StartServer(*address, *port, ".", username, password, *verbose)
}
