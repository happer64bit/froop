package server

import (
	_ "embed"
	"html/template"
)

//go:embed views/browser.html
var browserTemplate string

func LoadTemplate() (*template.Template, error) {
	return template.New("browser").Parse(browserTemplate)
}
