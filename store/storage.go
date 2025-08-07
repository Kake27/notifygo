package store

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateStore struct {
	templates map[string] *template.Template
}

func NewTemplateStore() (*TemplateStore, error) {
	tStore := &TemplateStore{templates: make(map[string]*template.Template)}

	templateDir := filepath.Join(".", "store", "templates")

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".tmpl" {
			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return err
			}
			tStore.templates[info.Name()] = tmpl
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tStore, nil
}


func (ts *TemplateStore) GetTemplate(name string) (*template.Template, error) {
	tmpl, ok := ts.templates[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}
	return tmpl, nil
}