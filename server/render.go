package server

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type FileInfo struct {
	Name string
	Link string
}

type DirectoryData struct {
	Path  string
	Files []FileInfo
}

func serveDirectory(w http.ResponseWriter, tmpl *template.Template, fullPath, relPath string, verbose bool) {
	names, err := ReadDirNames(fullPath)
	if err != nil {
		if verbose {
			fmt.Printf("Read error: %s\n", err)
		}
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	var files []FileInfo
	for _, name := range names {
		link := filepath.ToSlash(filepath.Join(relPath, name))
		if IsDir(filepath.Join(fullPath, name)) {
			link += "/"
		}
		files = append(files, FileInfo{Name: name, Link: link})
	}

	data := DirectoryData{Path: relPath, Files: files}
	if verbose {
		fmt.Printf("Rendering: %+v\n", data)
	}
	if err := tmpl.Execute(w, data); err != nil && verbose {
		fmt.Printf("Template render error: %s\n", err)
	}
}
