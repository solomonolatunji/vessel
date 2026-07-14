package models

type UserComposeFile struct {
	Services map[string]UserComposeService `yaml:"services"`
	Networks map[string]struct {
		External bool `yaml:"external"`
	} `yaml:"networks"`
}

type UserComposeService struct {
	Image       string            `yaml:"image"`
	Build       any               `yaml:"build"`
	Ports       []string          `yaml:"ports"`
	Environment map[string]string `yaml:"environment"`
	EnvFile     string            `yaml:"env_file"`
	Volumes     []string          `yaml:"volumes"`
	DependsOn   []string          `yaml:"depends_on"`
	Restart     string            `yaml:"restart"`
}
