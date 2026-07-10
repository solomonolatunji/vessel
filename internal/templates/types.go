package templates

// ComposeTemplate represents a standard Docker Compose structure.
type ComposeTemplate struct {
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]interface{}    `yaml:"volumes,omitempty"`
}

// ComposeService represents a single service definition in a docker-compose file.
type ComposeService struct {
	Image       string          `yaml:"image"`
	Environment []string        `yaml:"environment,omitempty"`
	Ports       []string        `yaml:"ports,omitempty"`
	Volumes     []string        `yaml:"volumes,omitempty"`
	Command     []string        `yaml:"command,omitempty"`
	DependsOn   []string        `yaml:"depends_on,omitempty"`
	XVessel     *VesselMetadata `yaml:"x-vessel,omitempty"`
}

// VesselMetadata contains Vessel-specific metadata for a service
type VesselMetadata struct {
	IsDatabase       bool                  `yaml:"is_database,omitempty"`
	ConnectionString string                `yaml:"connection_string,omitempty"`
	Backup           *VesselBackupMetadata `yaml:"backup,omitempty"`
}

// VesselBackupMetadata contains backup execution instructions
type VesselBackupMetadata struct {
	Command       []string `yaml:"command,omitempty"`
	FileExtension string   `yaml:"file_extension,omitempty"`
}
