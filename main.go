package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/akamensky/argparse"

	"net"

	"github.com/happer64bit/froop/server"
)

func smart_replace_address(address string) string {
	// If address is "localhost", treat it as empty to auto-detect LAN IP
	if address == "" || address == "localhost" {
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			var lanIP string
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					ip := ipnet.IP
					// Prefer LAN (private) addresses
					if ip.IsPrivate() {
						return ip.String()
					}
					// Fallback to first non-loopback IPv4 if no private found
					if lanIP == "" {
						lanIP = ip.String()
					}
				}
			}
			if lanIP != "" {
				return lanIP
			}
		}
		return "127.0.0.1"
	}
	return address
}

func main() {
	// Create a new parser object
	parser := argparse.NewParser("froop", "A Froop Http Cli")

	// Create the 'serve' command
	serveCmd := parser.NewCommand("serve", "Start the HTTP server")
	
	// Add command-line options for port and address
	port := serveCmd.String("p", "port", &argparse.Options{
		Help:    "Port to run the server on",
		Default: "8080",
	})

	address := serveCmd.String("a", "address", &argparse.Options{
		Help:    "Address to bind the server to",
		Required: false,
	})

	// Add an optional authentication flag
	auth := serveCmd.String("", "auth", &argparse.Options{
		Help:     "Enable authentication with the format username:password",
		Required: false,
	})

	// Add a verbose flag for logging
	verbose := serveCmd.Flag("v", "verbose", &argparse.Options{
		Help: "Enable verbose logging",
	})

	// Parse the command-line arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	// Check if the 'serve' command was provided
	if len(os.Args) < 2 || os.Args[1] != "serve" {
		fmt.Println("Usage: <program> serve [--port=PORT] [--address=ADDRESS] [--auth=username:password] [--verbose]")
		os.Exit(1)
	}

	// Handle authentication parsing
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

	smart_address := smart_replace_address(*address)

	server.StartServer(smart_address, *port, ".", username, password, *verbose)
}
