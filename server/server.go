package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mdp/qrterminal/v3"
)

type FileInfo struct {
	Name string
	Link string
}

type DirectoryData struct {
	Path  string
	Files []FileInfo
}

func BasicAuthMiddleware(next http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func StartServer(address, port, staticDir, authUsername, authPassword string) {
	addr := fmt.Sprintf("%s:%s", address, port)

	fmt.Printf("Generating QR code for: http://%s\n", addr)
	qrterminal.GenerateWithConfig("http://"+addr, qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	})

	tmpl, err := template.ParseFiles("./views/browser.html")
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Clean the URL path and prepend the static directory
		requestedPath := filepath.Join(staticDir, r.URL.Path)
		fmt.Printf("Requested path: %s\n", requestedPath)

		fileInfo, err := os.Stat(requestedPath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "File or directory not found", http.StatusNotFound)
			} else {
				fmt.Printf("Error stating path '%s': %s\n", requestedPath, err)
				http.Error(w, "Error accessing file or directory", http.StatusInternalServerError)
			}
			return
		}

		if fileInfo.IsDir() {
			// Serve directory listing
			dir, err := os.Open(requestedPath)
			if err != nil {
				fmt.Printf("Error opening directory '%s': %s\n", requestedPath, err)
				http.Error(w, "Error opening directory", http.StatusInternalServerError)
				return
			}
			defer func() {
				if cerr := dir.Close(); cerr != nil {
					fmt.Printf("Error closing directory '%s': %s\n", requestedPath, cerr)
				}
			}()

			fileNames, err := dir.Readdirnames(0)
			if err != nil {
				fmt.Printf("Error reading directory '%s': %s\n", requestedPath, err)
				http.Error(w, "Error reading directory", http.StatusInternalServerError)
				return
			}

			fmt.Printf("Files found in directory '%s': %v\n", requestedPath, fileNames)

			var files []FileInfo
			for _, fileName := range fileNames {
				filePath := filepath.Join(r.URL.Path, fileName)
				if fileInfo, err := os.Stat(filepath.Join(requestedPath, fileName)); err == nil && fileInfo.IsDir() {
					filePath += "/" // Append '/' for directories
				}
				files = append(files, FileInfo{
					Name: fileName,
					Link: filePath,
				})
			}

			data := DirectoryData{
				Path:  r.URL.Path,
				Files: files,
			}

			fmt.Printf("Rendering template with data: %+v\n", data)
			err = tmpl.Execute(w, data)
			if err != nil {
				fmt.Printf("Error rendering template: %s\n", err)
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
		} else {
			http.ServeFile(w, r, requestedPath)
		}
	})

	var handler http.Handler = http.DefaultServeMux
	if authUsername != "" && authPassword != "" {
		handler = BasicAuthMiddleware(handler, authUsername, authPassword)
	}

	fmt.Printf("Starting server at http://%s\n", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
