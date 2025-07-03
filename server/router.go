package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleRoutes(staticDir string, tmpl *template.Template, verbose bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		relPath := strings.TrimPrefix(r.URL.Path, "/")
		if relPath == "" {
			relPath = "."
		}

		relPath = filepath.ToSlash(relPath)
		fullPath := filepath.Join(staticDir, relPath)

		if verbose {
			fmt.Printf("Requested: %s\n", fullPath)
		}

		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
			} else {
				if verbose {
					fmt.Printf("Stat error: %s\n", err)
				}
				http.Error(w, "Error accessing file", http.StatusInternalServerError)
			}
			return
		}

		if info.IsDir() {
			serveDirectory(w, tmpl, fullPath, relPath, verbose)
		} else {
			http.ServeFile(w, r, fullPath)
		}
	})
}
