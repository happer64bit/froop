package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
)

// FileInfo holds the file information for the template
type FileInfo struct {
	Name string
	Link string
}

// DirectoryData holds the path and the list of files for the template
type DirectoryData struct {
	Path  string
	Files []FileInfo
}

// BasicAuthMiddleware is a middleware for Basic Authentication
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

func StartServer(address, port, staticDir, authUsername, authPassword string, verbose bool) {
	addr := fmt.Sprintf("%s:%s", address, port)

	qr_obj := qrcodeTerminal.New()

	qr_obj.Get("http://" + addr).Print()

	fmt.Printf("Generating QR code for: http://%s\n", addr)

	tmpl, err := template.ParseFiles("./views/browser.html")
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		return
	}

	// Serve the directory listing or file content
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Clean the URL path and prepend the static directory
		requestedPath := filepath.Join(staticDir, r.URL.Path)
		if verbose {
			fmt.Printf("Requested path: %s\n", requestedPath)
		}

		// Check if requested path is a directory
		fileInfo, err := os.Stat(requestedPath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "File or directory not found", http.StatusNotFound)
			} else {
				if verbose {
					fmt.Printf("Error stating path '%s': %s\n", requestedPath, err)
				}
				http.Error(w, "Error accessing file or directory", http.StatusInternalServerError)
			}
			return
		}

		if fileInfo.IsDir() {
			// Serve directory listing
			dir, err := os.Open(requestedPath)
			if err != nil {
				if verbose {
					fmt.Printf("Error opening directory '%s': %s\n", requestedPath, err)
				}
				http.Error(w, "Error opening directory", http.StatusInternalServerError)
				return
			}
			defer func() {
				if cerr := dir.Close(); cerr != nil {
					if verbose {
						fmt.Printf("Error closing directory '%s': %s\n", requestedPath, cerr)
					}
				}
			}()

			fileNames, err := dir.Readdirnames(0)
			if err != nil {
				if verbose {
					fmt.Printf("Error reading directory '%s': %s\n", requestedPath, err)
				}
				http.Error(w, "Error reading directory", http.StatusInternalServerError)
				return
			}

			if verbose {
				fmt.Printf("Files found in directory '%s': %v\n", requestedPath, fileNames)
			}

			// Prepare the list of files for the template
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

			// Create the data to pass to the template
			data := DirectoryData{
				Path:  r.URL.Path,
				Files: files,
			}

			// Render the template with the directory data
			if verbose {
				fmt.Printf("Rendering template with data: %+v\n", data)
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				if verbose {
					fmt.Printf("Error rendering template: %s\n", err)
				}
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
		} else {
			// Serve file content
			http.ServeFile(w, r, requestedPath)
		}
	})

	// Apply Basic Authentication middleware if credentials are provided
	var handler http.Handler = http.DefaultServeMux
	if authUsername != "" && authPassword != "" {
		handler = BasicAuthMiddleware(handler, authUsername, authPassword)
	}

	// Log the server start and listen on the given address and port
	if verbose {
		fmt.Printf("Starting server at http://%s\n", addr)
	}
	if err := http.ListenAndServe(addr, handler); err != nil {
		if verbose {
			fmt.Printf("Error starting server: %s\n", err)
		}
		os.Exit(1)
	}
}
