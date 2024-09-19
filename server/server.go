package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mdp/qrterminal/v3"
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

func StartServer(address, port, staticDir string) {
	addr := fmt.Sprintf("%s:%s", address, port)

	// Generate QR code for the server address
	fmt.Printf("Generating QR code for: http://%s\n", addr)
	qrterminal.GenerateWithConfig("http://"+addr, qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	})

	// Parse the HTML template
	tmpl, err := template.ParseFiles("./views/browser.html")
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		return
	}

	// Serve the directory listing or file content
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Clean the URL path and prepend the static directory
		requestedPath := filepath.Join(staticDir, r.URL.Path)
		fmt.Printf("Requested path: %s\n", requestedPath)

		// Check if requested path is a directory
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

			// Prepare the list of files for the template
			var files []FileInfo
			for _, fileName := range fileNames {
				filePath := filepath.Join(r.URL.Path, fileName) // Relative link for files
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
			fmt.Printf("Rendering template with data: %+v\n", data)
			err = tmpl.Execute(w, data)
			if err != nil {
				fmt.Printf("Error rendering template: %s\n", err)
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
			}
		} else {
			// Serve file content
			http.ServeFile(w, r, requestedPath)
		}
	})

	// Log the server start and listen on the given address and port
	fmt.Printf("Starting server at http://%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
