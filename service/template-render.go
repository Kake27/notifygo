package service

import (
	"bytes"
	"fmt"
	"text/template"
	"notification-service/store"
)

type TemplateRenderer struct {
	fileStore *store.TemplateStore
	dbStore *store.DBTemplateStore
}

func NewTemplateRenderer(file *store.TemplateStore, db *store.DBTemplateStore) *TemplateRenderer {
	return &TemplateRenderer{fileStore: file, dbStore: db}
}

func (r *TemplateRenderer) Render(templateName string, data map[string]string) (string, error) {
	if r.fileStore!=nil {
		if tmpl, err := r.fileStore.GetTemplate(templateName); err == nil {
			var buf bytes.Buffer
			_ = tmpl.Execute(&buf, data)
			return buf.String(), nil
		}
	}

	dbTemplate, err := r.dbStore.GetByName(templateName)
	if err!=nil {
		return "", fmt.Errorf("failed to get template from DB: %w", err)
	}

	tmplParsed, err := template.New("dbTemplate").Parse(dbTemplate.Content)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	var buf bytes.Buffer
	err = tmplParsed.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}