package templates

type ComposeTemplate struct {
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]interface{}    `yaml:"volumes,omitempty"`
}

type ComposeService struct {
	Image       string                 `yaml:"image"`
	Environment map[string]string      `yaml:"environment,omitempty"`
	Ports       []string               `yaml:"ports,omitempty"`
	Volumes     []string               `yaml:"volumes,omitempty"`
	Command     []string               `yaml:"command,omitempty"`
	DependsOn   []string               `yaml:"depends_on,omitempty"`
	Restart     string                 `yaml:"restart,omitempty"`
	HealthCheck map[string]interface{} `yaml:"healthcheck,omitempty"`
}
