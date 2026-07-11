package templates

type ComposeTemplate struct {
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]interface{}    `yaml:"volumes,omitempty"`
}

type ComposeService struct {
	Image       string          `yaml:"image"`
	Environment []string        `yaml:"environment,omitempty"`
	Ports       []string        `yaml:"ports,omitempty"`
	Volumes     []string        `yaml:"volumes,omitempty"`
	Command     []string        `yaml:"command,omitempty"`
	DependsOn   []string        `yaml:"depends_on,omitempty"`
	XVessel     *VesselMetadata `yaml:"x-vessel,omitempty"`
}

type VesselMetadata struct {
	IsDatabase       bool                  `yaml:"is_database,omitempty"`
	ConnectionString string                `yaml:"connection_string,omitempty"`
	Backup           *VesselBackupMetadata `yaml:"backup,omitempty"`
}

type VesselBackupMetadata struct {
	Command       []string `yaml:"command,omitempty"`
	FileExtension string   `yaml:"file_extension,omitempty"`
}
