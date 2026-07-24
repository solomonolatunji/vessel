package engine

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"

	"codedock.run/codedock/internal/utils"
)

//go:embed compose/*.yaml
var templateFiles embed.FS

type TemplateManager struct {
	templates map[string]ComposeTemplate
}

func NewTemplateManager() (*TemplateManager, error) {
	mgr := &TemplateManager{
		templates: make(map[string]ComposeTemplate),
	}

	err := fs.WalkDir(templateFiles, "compose", mgr.walkDir)
	if err != nil {
		fmt.Printf("warning: template manager WalkDir error: %v\n", err)
	}

	return mgr, nil
}

func (m *TemplateManager) walkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
		return nil
	}

	data, err := templateFiles.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", path, err)
	}

	var tmpl ComposeTemplate
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		fmt.Printf("warning: failed to parse template %s: %v\n", path, err)
		return nil // skip bad files instead of failing everything
	}

	id := strings.TrimSuffix(d.Name(), ".yaml")
	m.templates[id] = tmpl
	return nil
}

func (m *TemplateManager) GetTemplate(id string) (ComposeTemplate, error) {
	if m == nil || m.templates == nil {
		return ComposeTemplate{}, utils.NewNotFoundError("Template", id)
	}
	tmpl, exists := m.templates[id]
	if !exists {
		return ComposeTemplate{}, utils.NewNotFoundError("Template", id)
	}
	return tmpl, nil
}

func (m *TemplateManager) ListTemplates() []string {
	if m == nil || m.templates == nil {
		return []string{}
	}
	list := make([]string, 0, len(m.templates))
	for id := range m.templates {
		list = append(list, id)
	}
	return list
}
