package engine

type ComposeTemplate struct {
	Version   string                    `yaml:"version,omitempty"`
	Services  map[string]ComposeService `yaml:"services"`
	Volumes   map[string]interface{}    `yaml:"volumes,omitempty"`
	XCodedock *CodedockMetadata         `yaml:"x-codedock,omitempty"`
}

type ComposeService struct {
	Image       string            `yaml:"image"`
	Environment []string          `yaml:"environment,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Command     []string          `yaml:"command,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	XCodedock   *CodedockMetadata `yaml:"x-codedock,omitempty"`
}

type CodedockMetadata struct {
	IsDatabase       bool                     `yaml:"is_database,omitempty"`
	IsOneClick       bool                     `yaml:"is_one_click,omitempty"`
	Name             string                   `yaml:"name,omitempty"`
	Description      string                   `yaml:"description,omitempty"`
	ConnectionString string                   `yaml:"connection_string,omitempty"`
	Backup           *CodedockBackupMetadata  `yaml:"backup,omitempty"`
	Restore          *CodedockRestoreMetadata `yaml:"restore,omitempty"`
}

type CodedockBackupMetadata struct {
	Command       []string `yaml:"command,omitempty"`
	FileExtension string   `yaml:"file_extension,omitempty"`
}

type CodedockRestoreMetadata struct {
	Command []string `yaml:"command,omitempty"`
}
