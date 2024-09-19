package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

// GetExecutablePath returns the directory of the currently running executable
func GetExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

func StartServer(address, port, staticDir, authUsername, authPassword string, verbose bool) {
	addr := fmt.Sprintf("%s:%s", address, port)

	qr_obj := qrcodeTerminal.New()
	qr_obj.Get("http://" + addr).Print()

	fmt.Printf("Starting server at http://%s\n", addr)

	binaryDir, err := GetExecutablePath()
	if err != nil {
		fmt.Printf("Error getting executable path: %s\n", err)
		return
	}

	// Build the path to the template
	viewsPath := filepath.Join(binaryDir, "views", "browser.html")

	// Parse the template using the absolute path
	tmpl, err := template.ParseFiles(viewsPath)
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		return
	}

	// Serve the directory listing or file content
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Sanitize and normalize the URL path
		urlPath := strings.TrimPrefix(r.URL.Path, "/")
		if urlPath == "" {
			urlPath = "."
		}

		// Normalize URL path to use forward slashes and avoid encoding issues
		urlPath = filepath.ToSlash(urlPath)
		requestedPath := filepath.Join(staticDir, urlPath)

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
				filePath := filepath.ToSlash(filepath.Join(urlPath, fileName))
				if fileInfo, err := os.Stat(filepath.Join(requestedPath, fileName)); err == nil && fileInfo.IsDir() {
					filePath += "/" // Append '/' for directories
				}
				files = append(files, FileInfo{
					Name: fileName,
					Link: filePath,
				})
			}

			data := DirectoryData{
				Path:  urlPath,
				Files: files,
			}

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

	if err := http.ListenAndServe(addr, handler); err != nil {
		if verbose {
			fmt.Printf("Error starting server: %s\n", err)
		}
		os.Exit(1)
	}
}
