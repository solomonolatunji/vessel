package types

// ProjectCanvasSummary aggregates all services, environments, and online counts for a Project card on the dashboard.
type ProjectCanvasSummary struct {
	ProjectConfig
	EnvironmentsCount  int                `json:"environmentsCount"`
	AppsCount          int                `json:"appsCount"`
	DatabasesCount     int                `json:"databasesCount"`
	StorageCount       int                `json:"storageCount"`
	OnlineServices     int                `json:"onlineServices"`
	TotalServices      int                `json:"totalServices"`
	ServiceIcons       []string           `json:"serviceIcons"` // e.g. ["github", "postgres", "redis"]
	DefaultEnvironment *EnvironmentConfig `json:"defaultEnvironment,omitempty"`
}

// EnvironmentCanvas aggregates all application services, databases, and storage buckets belonging to a specific environment.
type EnvironmentCanvas struct {
	Environment *EnvironmentConfig  `json:"environment"`
	Apps        []*AppServiceConfig `json:"apps"`
	Databases   []*DatabaseConfig   `json:"databases"`
	Storage     []*StorageConfig    `json:"storage"`
}
