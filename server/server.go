package server

import (
	"fmt"
	"net/http"
	"os"
)

func StartServer(address, port, staticDir, authUser, authPass string, verbose bool) {
	addr := fmt.Sprintf("%s:%s", address, port)
	PrintQRCode(addr)

	tmpl, err := LoadTemplate()
	if err != nil {
		fmt.Printf("Template error: %s\n", err)
		return
	}

	HandleRoutes(staticDir, tmpl, verbose)

	var handler http.Handler = http.DefaultServeMux
	if authUser != "" && authPass != "" {
		handler = BasicAuthMiddleware(handler, authUser, authPass)
	}

	if err := http.ListenAndServe(addr, handler); err != nil {
		if verbose {
			fmt.Printf("Server error: %s\n", err)
		}
		os.Exit(1)
	}
}
