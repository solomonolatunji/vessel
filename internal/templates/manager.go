package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

var templateFiles embed.FS

type Manager struct {
	templates map[string]ComposeTemplate
}

func NewManager() (*Manager, error) {
	mgr := &Manager{
		templates: make(map[string]ComposeTemplate),
	}

	err := fs.WalkDir(templateFiles, "compose", mgr.walkDir)
	if err != nil {
		return nil, err
	}

	return mgr, nil
}

func (m *Manager) walkDir(path string, d fs.DirEntry, err error) error {
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
		return fmt.Errorf("failed to parse template %s: %w", path, err)
	}

	id := strings.TrimSuffix(d.Name(), ".yaml")
	m.templates[id] = tmpl
	return nil
}

func (m *Manager) GetTemplate(id string) (ComposeTemplate, error) {
	tmpl, exists := m.templates[id]
	if !exists {
		return ComposeTemplate{}, fmt.Errorf("template not found: %s", id)
	}
	return tmpl, nil
}

func (m *Manager) ListTemplates() []string {
	list := make([]string, 0, len(m.templates))
	for id := range m.templates {
		list = append(list, id)
	}
	return list
}
